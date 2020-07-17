# firestore-rules

The goal of this project is to write firestore rules in an enhanced language that specifies the types of docs. 
Then run a compiler that reads the enhanced language and writes standard "filestores.rules" that can be installed in firestore.

Cloud firestore rules files are kind of hard to work with. 
First, the only way to check a rule file is to install it in firestore, and firestore error messages are not terribly friendly.
Second, there is no type checking until runtime.
Third, writing rules to disallow malicious writes is tedious. 
Every object that is written should have at least a size test to make sure a malicious user doesn't dump a lot of junk into docs.
For even a fairly simple doc, this can quickly become a lot of code. 

The first milestone is to write a parser for firestore.rules that generates decent error messages.
I have a good start there, although I'm sure there are some bugs.

The second milestone is to add syntax to specify the type of data on each page. 
This amounts to specifying what keys are allowed in each doc, and the types and maximum sizes of the associated values.
I haven't worked out the exact syntax for this yet, but it will probably look like a very restricted version of a typescript type declaration.

The type specification will be used to generate validation code that is added to "allow" statements to prevent bad data from being written.
E.g., the compiler could add to each "match" statement a dataIsValid() function that checks request.resource.data for validity. A call to
the validation function can be automatically added to the "allow create: if ..." and "allow update: if ...". 













