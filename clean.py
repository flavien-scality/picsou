import subprocess

instances = subprocess.check_output("aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId'", shell=True)
instance_ids = []
for instance in instances.split('"'):
  if "i" in instance:
    instance_ids.append(instance)
for iid in instance_ids:
  subprocess.call("aws ec2 terminate-instances --instance-ids {}".format(iid), shell=True)
