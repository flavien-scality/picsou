"""
Lambda to kill EBS orphans on all AWS regions
"""

import boto3

import logging

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def get_linked_sgs():
    """
    Iterator to get security groups linked to instances

    :return: linked security groups ids
    :trype: Iterable(string)
    """
    ec2 = boto3.client("ec2")
    instances_list = ec2.describe_instances(DryRun=False)
    for reservation in instances_list["Reservations"]:
        for instance in reservation["Instances"]:
            for sg in instance["SecurityGroups"]:
                yield sg["GroupId"]
    return



def get_orphan_sgs():
    """
    Iterator to get orphan security groups

    :return: orphan security groups
    :rtype: Iterable(boto3.ec2.SecurityGroup)
    """
    ec2 = boto3.client("ec2")
    ec2_r = boto3.resource("ec2")
    sgs = ec2.describe_security_groups(DryRun=False)
    for sg in sgs["SecurityGroups"]:
        if any(group_id not in sg["GroupId"] for group_id in list(get_linked_sgs())):
            yield ec2_r.SecurityGroup(sg["GroupId"])
    return



def handle(event, context):
    """
    Lambda handler
    """
    logger.info("event: {} | context {}".format(event, context))
    count = 1
    for orphan in get_orphan_sgs():
        logger.info("Orphan sg {} to delete: {}".format(count, orphan))
        try:
            orphan.delete()
            logger.info("Orphan sg {} deleted".format(count))
            count += 1
        except Exception as e:
            logger.warning("Could not delete orphan sg: {}".format(orphan.group_id))
    logger.info("number of orphan volumes deleted: {}".format(count))
