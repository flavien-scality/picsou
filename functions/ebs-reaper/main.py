"""
Lambda to kill EBS orphans on all AWS regions
"""

import boto3

import logging

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def get_orphan_volumes():
    """
    Iterator to get orphan volumes

    :return: orphan volumes
    :rtype: Iterable(boto3.ec2.Volume)
    """
    ec2 = boto3.client("ec2")
    volumes = ec2.describe_volumes(DryRun=False)
    for volume in volumes["Volumes"]:
        logger.debug("volume: {}".format(volume))
        if len(volume["Attachments"]) == 0:
            yield volume
    return



def handle(event, context):
    """
    Lambda handler
    """
    ec2 = boto3.resource("ec2")
    logger.info("event: {} | context {}".format(event, context))
    count = 0
    for orphan in get_orphan_volumes():
        logger.info("Orphan {} to delete: {}".format(count, orphan))
        ec2.Volume(orphan["VolumeId"]).delete()
        logger.info("orphan {} deleted".format(count))
        count += 1
    logger.info("number of orphan volumes deleted: {}".format(count))
