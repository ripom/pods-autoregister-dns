FROM python:3.7.0-stretch
RUN pip install azure-mgmt-dns
RUN pip install azure-mgmt-privatedns
WORKDIR /code
copy dnsregister.py dnsregister.py
CMD ["python","/code/dnsregister.py"]
