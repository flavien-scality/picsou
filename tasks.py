#!/usr/bin/env python3

from invoke import task
from pprint import pprint

import os

def terraform_init(ctx):
  ctx.run("terraform init -force-copy -backend-config=vars.tfvars")
  ctx.run("terraform state pull")

def get_lambda_paths(ctx):
  output = ctx.run("ls functions")
  for lambda_name in output.stdout.split("\n"):
    if len(lambda_name) == 0:
      continue
    script_dir = "{0}/functions/{1}".format(os.getcwd(), lambda_name)
    yield (lambda_name, script_dir)

@task
def deploy_lambdas(ctx, state="test"):
  for (lambda_name, lambda_path) in get_lambda_paths(ctx):
    with ctx.cd(lambda_path):
      if os.path.exists("{0}.zip".format(lambda_name)):
        ctx.run("mv {0}.zip {0}.zip.backup".format(lambda_name))
      venv_path = "{0}/venv".format(lambda_path)
      if os.path.exists(venv_path):
        ctx.run("rm -rf {0}".format(venv_path))
      ctx.run("python3.6 -m venv venv")
      with ctx.prefix("source venv/bin/activate"):
        ctx.run("pip install -r requirements.txt")
      ctx.run("zip -r9 {0}/{1}.zip {0}/venv/lib/python3.6/site-packages/*".format(lambda_path, lambda_name))
      ctx.run("zip -r9 {0}/{1}.zip {0}/{1}.py".format(lambda_path, lambda_name))
      terraform_init(ctx)
      if state == "test":
        ctx.run("terraform plan")
      elif state == "deploy":
        ctx.run("terraform apply")


@task
def clean(ctx):
  for (lambda_name, lambda_path) in get_lambda_paths(ctx):
    with ctx.cd(lambda_path):
      terraform_init(ctx)
      ctx.run("terraform destroy -force")


@task
def deploy(ctx):
    deploy_lambdas(ctx, state="deploy")


@task
def test(ctx):
    deploy_lambdas(ctx, state="test")
