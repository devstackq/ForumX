    #!/bin/bash
    docker build --tag forumx .
    docker run --publish 8080:6969 --detach --name forum forumx