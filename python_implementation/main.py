import os
import io
import json
import xlsxwriter
import requests
import calendar
import time
from modules.s3service import S3Service

def handler(event, context):
    s3 = S3Service(
        aws_access_key_id=os.environ['aws_access_key_id'],
        aws_secret_access_key=os.environ["aws_secret_access_key"]
    )
    bucket = os.environ['bucket']
    upload_bucket = os.environ['upload_bucket']
    key = event['queryStringParameters']['file_name']
  
    file = s3.get_object(bucket, key)
    my_list = json.loads(file['Body'].read().decode('utf-8'))
    del file

    workbook = xlsxwriter.Workbook('/tmp/test.xlsx', {'constant_memory': True})
    worksheet = workbook.add_worksheet()
    for col_num, col_data in enumerate(my_list[0]):
        worksheet.set_column(col_num, 0, 35)
    for row_num, row_data in enumerate(my_list):
        for col_num, col_data in enumerate(row_data):
            worksheet.write(row_num, col_num, col_data)
    del my_list
    workbook.close()
    
    directory = event['queryStringParameters']['user_id']
    filename = str(calendar.timegm(time.gmtime())) + '.xlsx'
    report_id = event['queryStringParameters']['report_id']

    with open('/tmp/test.xlsx', 'rb') as f:
        s3.put_object(upload_bucket, f, directory + '/' + filename)
    
    if os.path.exists('/tmp/test.xlsx'):
        os.remove('/tmp/test.xlsx')
    endpoint = os.environ['endpoint']
    headers = {'Authorization': 'Bearer ' + os.environ['auth_token']}
    data = {'report_id': report_id, 'cloud_key': directory + '/' + filename}
    response = requests.post(endpoint, data=data, headers=headers).json()
    
    return {
        'statusCode': 200,
        'headers': {
            'Content-Type': 'text/plain'
        },
        'body': response
    }