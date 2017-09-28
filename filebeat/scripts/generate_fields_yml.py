import json
import re
import yaml

from argparse import ArgumentParser
from collections import OrderedDict


FIELDS_YML_PATH = "module/{module}/{fileset}/_meta/fields.yml"
PIPELINE_JSON_PATH = "module/{module}/{fileset}/ingest/pipeline.json"
DEFAULT_TYPE = "text"
TYPES = {
    "WORD": "keyword",
    "USERNAME": "keyword",
    "HOSTNAME": "keyword",
    "LOGLEVEL": "keyword",
    "GREEDYDATA": "text",
    "GREEDYMULTILINE": "text",
    "IPHOST": "keyword",
    "IPORHOST": "keyword",
    "NUMBER": "long",
}
FUNCTIONS = {
    "remove": None,
    "rename": None,
}


def _get_elements_of_grok_patterns(pipeline):
    elements = []
    patterns = pipeline["processors"][0]["grok"]["patterns"]

    for pattern in patterns:
        fields = re.findall("{[\.\w\:]*}", pattern)
        for field in fields:
            field_type, name = field[1:-1].split(":")
            elements.append((field_type, name))

    return elements


def get_message_elements_from_pipeline(module_name, fileset_name):
    pipeline_path = PIPELINE_JSON_PATH.format(module=module_name, fileset=fileset_name)
    with open(pipeline_path, "r") as pipeline_file:
        pipeline = json.load(pipeline_file, object_pairs_hook=OrderedDict)

    return _get_elements_of_grok_patterns(pipeline)


def __get_field_by_name(fields, field_name):
    for field in fields:
        if "name" in field and field["name"] == field_name:
            return field
    return None


def __insert_last_field(fields, element_type, element_name):
    field = __get_field_by_name(fields, element_name)
    if field:
        return

    field_type = DEFAULT_TYPE
    if element_type in TYPES:
        field_type = TYPES[element_type]
    new_field = {
                    "name": element_name,
                    "type": field_type,
                    "description": "Please add description",
                    "example": "Please provide an example",
                }
    fields.append(new_field)


def __insert_into_existing_group(field, field_elements, element_type, index, count):
    _insert_fields_into_dict(field["fields"], field_elements, element_type, index+1, count)


def __insert_new_group(fields, field_elements, element_type, index, count):
    its_fields = []
    _insert_fields_into_dict(its_fields, field_elements, element_type, index+1, count)

    if its_fields:
        fields.append({
            "name": field_elements[index],
            "type": "group",
            "description": "Please add description",
            "fields": its_fields,
        })


def __insert_group(fields, field_elements, element_type, index, count):
    field = __get_field_by_name(fields, field_elements[index])
    if field is not None:
        __insert_into_existing_group(field, field_elements, element_type, index, count)
    else:
        __insert_new_group(fields, field_elements, element_type, index, count)


def _insert_fields_into_dict(fields, field_elements, element_type, index, count):
    if index+1 == count:
        __insert_last_field(fields, element_type, field_elements[index])
        return
    __insert_group(fields, field_elements, element_type, index, count)


def get_nested_fields(message_elements):
    fields = []
    for element_type, name in message_elements:
        name_elements = name.split(".")
        elements_count = len(name_elements)
        _insert_fields_into_dict(fields, name_elements, element_type, 1, elements_count)

    return fields


def save_yml(fields_yml, module_name, fileset_name):
    fields_yml_path = FIELDS_YML_PATH.format(module=module_name, fileset=fileset_name)
    with open(fields_yml_path, "w") as out:
        yaml.safe_dump(fields_yml, out, default_flow_style=False, encoding="utf-8", allow_unicode=True, width=100)


def generate_fields_yml(module_name, fileset_name):
    message_elements = get_message_elements_from_pipeline(module_name, fileset_name)
    fields_yml = get_nested_fields(message_elements)
    save_yml(fields_yml, module_name, fileset_name)


if __name__ == "__main__":
    parser = ArgumentParser(description="Generate fields.yml for new Filebeat modules from pipeline.json")
    parser.add_argument("module")
    parser.add_argument("fileset")
    args = parser.parse_args()

    generate_fields_yml(args.module, args.fileset)
