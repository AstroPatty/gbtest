import boto3
import os
from datetime import datetime


def handler(event, context):
    TABLE_NAME = os.environ["DB_NAME"]
    dbResource = boto3.resource('dynamodb')
    table = dbResource.Table(TABLE_NAME)
    project_name = event["project_name"]
    response = table.get_item(TABLE_NAME, {"project-name": project_name})
    if "Item" in response:
        return {
            "code": 409,
            "msg": f"A project named {project_name} already exists."
        }

    table.put_item(
        Item = {
            "project-name": project_name,
            "created-data": datetime.utcnow().isoformat()

        }
    )

