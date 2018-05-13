========
toybox
========

A minimalist toolbox. (Respect busybox)

currently, only suports linux.

Available commands
===================

.. csv-table:: Status
   :header: "Command Name", "Status", "Test", "Desc"

   arp, completed, no
   base64, completed, no
   cat, completed, no
   chgrp, completed, no
   chmod, completed, no
   chown, completed, no
   cksum, bug, no
   cmp, completed, no
   cp, completed, no
   cut, not yet, yes, should arg parse manually (-f2 and so on)
   date, not yet, yes, `+` format is not implemented.
   diff, not yet, yes, diff algorithm is different from actual perhaps
   du, bug, yes, `-s` not work. size is incorrect(block size matter)
   dirname, completed, yes
   echo, completed, no
   false, completed, yes
   ls, not yet, no
   ln, completed, no
   mkdir, completed, no
   mv, completed, no
   md5sum, completed, no
   sha1sum, completed, no
   sha256sum, completed, no
   sha512sum, completed, no
   true, completed, yes
   rm, completed, no
   rmdir, completed, no
   which, completed, no
   wc, completed, no
   wget, completed, no
   yes, completed, no
   sleep, completed, no
   seq, completed, no

Not Portable
-----------------

- chgrp
- chown

TODO
=======

many

Memo
-----------

Single UNIX Specification v3, Shell and Utilities
http://www.unix.org/version3/xcu_contents.html
http://pubs.opengroup.org/onlinepubs/9699919799/idx/utilities.html

LICENSE
===================

Since this project uses many of Golang standard libraries, I choose the same License.

BSD 3-clause
