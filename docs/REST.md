# REST

* `/rest` umbrella group
    * `/doc` document headers
        * `GET` list all documents
        * `POST` create a document entry
        * `/{id}`
            * `GET` get a document header by id
            * `DELETE` delete a document entry and the content
                * `/content` document contents
                    * `GET` get a document content by id
                    * `POST` upload content