# firestore-rules

The goal of this project is to write firestore rules in an enhanced language that specifies the type of each doc
in the database so that basic data validation can be done automatically. 

When you create a cloud firestore database and open it to users, the only thing protecting you from a database full of
junk is your database rules. But writing rules to validate every create and update is tedious:
at the minimum, every object that is written should be checked for type and size. 
For even a fairly simple doc, this can quickly become a lot of code. 

The first milestone is to write a parser for firestore.rules that generates decent error messages.
I have a good start there, although I'm sure there are some bugs.

The second milestone is to add syntax to specify the type of data on each page. 
This amounts to specifying what keys are allowed in each doc, and the types and maximum sizes of the associated values.
I haven't worked out the exact syntax for this yet, but it will probably look like a very restricted version of a typescript type declaration, 
but with the addition of size restrictions.

The type specification will be used to generate validation code that is added to "allow" statements to prevent bad data from being written.
E.g., the compiler will add to each "match" statement a dataIsValid() function that checks request.resource.data for validity. A call to
the validation function can be automatically added to the "allow create: if ..." and "allow update: if ...".














