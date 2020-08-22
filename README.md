# xq

A simple XML parser that exposes the https://github.com/antchfx/xmlquery to the 
command line.  Use XPath queries to search for information in XML documents,
or execute basic XPath expressions to convert data.

## Find

To search for elements, use the `--find` command:

    cat sample.xml | ./xq --find '//MyElement[@someattr="abc"]'
    
## Exec

To execute an XPath expression, use the `--exec` command:

    cat sample.xml | ./xq --exec 'count(//MyElement)'
    cat sample.xml | ./xq --exec 'sum(//MyElement@bytes 
   
