# eps-conduit
HTTP/HTTPS Load Balancer written in GO

##Usage
By default, eps-conduit will bind to port 8000, but any port can be specified.
###Flags
* -b    list of backend hosts
  * ex:  eps-conduit -b "10.2.8.1, 10.2.8.2"
* -bind specify what port to bind to (defaults to 8000)
  * ex:  eps-conduit -bind 80
* -mode specifies what mode to use (http or https)
  * ex:  eps-conduit -mode https
* -cert specify an SSL cert file (for https mode)
  * ex:  eps-conduit -cert mycert.crt
* -key  specify an SSL keyfile (for https mode)
  * ex:  eps-conduit -key mykeyfile.key
