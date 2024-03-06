# SNMPv3 Enggine ID specification

```md
An SNMP engine's administratively-unique identifier.
Objects of this type are for identification, not for addressing, even though it is possible that an address may have been used in the generation of a specific value.
The value for this object may not be all zeros or all `ff`H or the empty (zero length) string.

The initial value for this object may be configured via an operator console entry or via an algorithmic function.
In the latter case, the following example algorithm is recommended.

In cases where there are multiple engines on the same system, the use of this algorithm is NOT appropriate, as it would result in all of those engines ending up with the same ID value.

1) The very first bit is used to indicate how the
  rest of the data is composed.

  0 - as defined by enterprise using former methods that existed before SNMPv3. See item 2 below.
  1 - as defined by this architecture, see item 3 below.

  Note that this allows existing uses of the engineID (also known as AgentID [RFC1910]) to co-exist with any new uses.

2) The snmpEngineID has a length of 12 octets.

  The first four octets are set to the binary equivalent of the agent's SNMP management private enterprise number as assigned by the Internet Assigned Numbers Authority (IANA).
  For example, if Acme Networks has been assigned { enterprises 696 }, the first four octets would be assigned `000002b8`H.

  The remaining eight octets are determined via one or more enterprise-specific methods.
  Such methods must be designed so as to maximize the possibility that the value of this object will be unique in the agent's administrative domain.
  For example, it may be the IP address of the SNMP entity, or the MAC address of one of the interfaces, with each address suitably padded with random octets.
  If multiple methods are defined, then it is recommended that the first octet indicate the method being used and the remaining octets be a function of the method.

3) The length of the octet string varies.

  The first four octets are set to the binary equivalent of the agent's SNMP management private enterprise number as assigned by the Internet Assigned Numbers Authority (IANA).
  For example, if Acme Networks has been assigned { enterprises 696 }, the first four octets would be assigned `000002b8`H.

  The very first bit is set to 1.
  For example, the above value for Acme Networks now changes to be `800002b8`H.

  The fifth octet indicates how the rest (6th and following octets) are formatted.
  The values for the fifth octet are:

    0     - reserved, unused.
    1     - IPv4 address (4 octets)
            lowest non-special IP address
    2     - IPv6 address (16 octets)
            lowest non-special IP address
    3     - MAC address (6 octets)
            lowest IEEE MAC address, canonical
            order
    4     - Text, administratively assigned
            Maximum remaining length 27
    5     - Octets, administratively assigned
            Maximum remaining length 27
    6-127 - reserved, unused
  128-255 - as defined by the enterprise
            Maximum remaining length 27
```
