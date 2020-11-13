    #!/bin/bash
    docker build --tag forumx .
    docker run --publish 8080:8080 --detach --name forum forumx