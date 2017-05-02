"""
Lambda to kill EBS victims on all AWS regions
"""

import boto3

import datetime
import json
import logging
import time

logging.basicConfig()
logger = logging.getLogger()
logger.setLevel(logging.INFO)

ZERO = datetime.timedelta(0)

class UTC(datetime.tzinfo):
  def utcoffset(self, dt):
    return ZERO
  def tzname(self, dt):
    return "UTC"
  def dst(self, dt):
    return ZERO

def datetime_handler(x):
    if isinstance(x, datetime.datetime):
        return x.isoformat()
    raise TypeError("Unknown type")

def get_instances():
    """
    Iterator to get ec2 instances from reservations list

    :param reservations: dict of ec2 reservations
    :type reservations: dict
    :return: ec2 instances id
    :rtype: Iterable(string)
    """
    ec2 = boto3.client("ec2")
    instances_list = ec2.describe_instances(DryRun=False)
    for reservation in instances_list["Reservations"]:
        for instance in reservation["Instances"]:
                logger.debug("victim metadatas: {}".format(json.dumps(instance, indent=4, sort_keys=True, default=datetime_handler)))
                yield instance["InstanceId"]
    return



def get_victims():
    """
    Iterator to get ec2 instances to kill

    :return: ec2 instances to kill
    :rtype: Iterable(boto3.ec2.Instance)
    """
    ec2 = boto3.client("ec2")
    ec2_r = boto3.resource("ec2")
    for instance in get_instances():
        logger.debug("instance: {}".format(instance))
        yield ec2_r.Instance(instance)
    return

def main():
    """
    Main function
    """
    count = 0
    err = 0
    for victim in get_victims():
        logger.info("Victim {} to delete: {}".format(count, victim))
        uptime = datetime.datetime.now(UTC()) - victim.launch_time
        try:
            # victim.delete()
            logger.info("instance launch deltatime: {}".format(uptime))
            logger.info("Victim {} deleted".format(count))
            count += 1
        except Exception as e:
            logger.warning("Err {}: Could not delete victim: {}".format(err, victim.instance_id))
            err += 1
    logger.info("number of victim instances deleted: {}, number of victim instances which could not be deleted: {}".format(count, err))

def handle(event, context):
    """
    Lambda handler
    """
    logger.info("event: {} | context {}".format(event, context))
    main()

if __name__ == "__main__":
    main()
