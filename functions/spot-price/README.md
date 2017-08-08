# Spot Instance Lambda

## Usage

### Shell example

```bash
aws lambda invoke \
--invocation-type RequestResponse \
--function-name spot-price \
--region eu-west-2 \
--log-type Tail \
--payload '{"InstanceTypes": [ "c4.xlarge" ]}' \
--profile default \
sport-price-output.txt
```

### Python example

```python
import boto3
```
