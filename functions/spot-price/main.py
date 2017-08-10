"""
Lambda to get best price for an ec2 instance filling specific conditions
"""

import boto3

from datetime import datetime, timedelta
import json
import logging
import re
import time
import urllib2

logging.basicConfig()
logger = logging.getLogger()
log_level = os.getenv("LOG_LEVEL", logging.INFO)
logger.setLevel(log_level)

def datetime_handler(x):
    if isinstance(x, datetime.datetime):
        return x.isoformat()
    raise TypeError("Unknown type")

def get_price(DryRun=True, event=None):
    """
    Get last price for different metadatas given

    :param DryRun: dry run aws api calls
    :type DryRun: boolean
    :param event: event given to the lambda including metadatas
                  to get efficient price
    :type event: dict
    :return: prices matching event metadatas
    :rtype: Iterable(float)
    """
    url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/eu-west-2/index.json"
    req = urllib2.Request(url)
    opener = urllib2.build_opener()
    f = opener.open(req)
    prices = json.loads(f.read())

    client = boto3.client("ec2")
    response = client.describe_spot_price_history(
        DryRun=DryRun,
        StartTime=event.get("StartTime", ""),
        EndTime=event.get("EndTime", ""),
        InstanceTypes=[ "|".join(event.get("InstanceTypes", "")) ],
        ProductDescriptions=event.get("ProductDescriptions", []),
        Filters=event.get("Filters", []),
        AvailabilityZone=event.get("AvailabilityZone", ""),
        MaxResults=event.get("MaxResults", 1000),
        NextToken=event.get("NextToken", "")
    )
    sprice = 1.0
    for sprices in response["SpotPriceHistory"]:
        if float(sprices["SpotPrice"]) < sprice:
            sprice = float(sprices["SpotPrice"])
    for content in prices["terms"]["OnDemand"]:
        for term in prices["terms"]["OnDemand"][content]:
            for price in prices["terms"]["OnDemand"][content][term]["priceDimensions"]:
                metadatas = prices["terms"]["OnDemand"][content][term]["priceDimensions"][price]["description"]
                # TODO: check for different tenancies (ex: dedicated)
                if "Linux" in metadatas and "Host" not in metadatas and "On Demand" in metadatas:
                    instance = re.findall("\w+\.\w+", metadatas)
                    price = float(instance[0])
                    instance_type = instance[1]
                    # instance_tenancy = re.findall("(On Demand|Dedicated)", metadatas)[0]
                    if any(instance_type in itype for itype in event.get("InstanceTypes", [])):
                        result = sprice + 0.05 * price
    return {
        "spotPrice": sprice,
        "onDemandPrice": price,
        "advisePrice": result
    }
    # regex to get [price, instances type]: ([^\s]+\.[^\s]+)
    # regex to get tenancy: (On Demand|Dedicated)
    # print("\n\nPrices: {}".format(json.dumps(prices["terms"], indent=4)))
    # print("\n\nPrices: {}".format(prices["terms"]["OnDemand"].keys()))


def handle(event, context):
    """
    Lambda handler
    """
    # logger.info("event: {} | context {}".format(event, context))
    return get_price(DryRun=False, event=event)


if __name__ == "__main__":
    today = datetime.now()
    event = {
         "Filters": [
            {
                'Name': 'product-description',
                'Values': [
                    'Linux/UNIX',
                ]
            }
        ],
        "StartTime": today - timedelta(days=1),
        "StartTime": today,
        "EndTime": today,
        "InstanceTypes": [ "c4.xlarge" ],
    }
    get_price(DryRun=False, event=event)
