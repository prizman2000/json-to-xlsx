import boto3


class S3Service(object):
    s3 = {}
    aws_access_key_id = ''
    aws_secret_access_key = ''

    # init singleton
    def __new__(cls, aws_access_key_id='', aws_secret_access_key=''):
        if not hasattr(cls, 'instance'):
            cls.instance = super(S3Service, cls).__new__(cls)
        if aws_access_key_id and aws_secret_access_key:
            cls.aws_access_key_id = aws_access_key_id
            cls.aws_secret_access_key = aws_secret_access_key
            cls.s3 = {}
        if not cls.s3:
            cls.create_session(cls)
        return cls.instance

    def create_session(self):
        session = boto3.session.Session(
            aws_access_key_id=self.aws_access_key_id,
            aws_secret_access_key=self.aws_secret_access_key,
            region_name='ru-central1'
        )

        self.s3 = session.client(
            service_name='s3',
            endpoint_url='https://storage.yandexcloud.net'
        )

        return self

    def put_object(self, bucket, body, key):
        return self.s3.put_object(
            Bucket=bucket,
            Body=body,
            Key=key
        )

    def get_object(self, bucket, key):
        return self.s3.get_object(
            Bucket=bucket,
            Key=key
        )

    def list_buckets(self):
        return self.s3.list_buckets()['Buckets']
