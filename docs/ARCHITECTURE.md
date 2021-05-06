# Architecture

The system's components are distributed.
Document _data_ is stored in an **S3** instance.
The _document index_, _permissions_ and _users_ are stored in an **AWS Aurora** instance.
The web page is hosted by the same server that holds the business logic, which runs in an **EC2** instance.
An **AWS Redis** instance is used as _session storage_.

