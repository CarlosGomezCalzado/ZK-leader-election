ZK-leader-election
===
# Introduction
This project implements a distributed leader election algorithm based on the proposal by [Carlos Gómez-Calzado et al.](https://addi.ehu.es/bitstream/handle/10810/16424/TESIS_CARLOS_GOMEZ_CALZADO.pdf?sequence=1). The algorithm utilizes two types of messages: a "birth" message, broadcasted using the reliable broadcast primitive when a process starts, and a leader election criterion based on the size of the "born" list and, in case of a tie, the process identifier.

# Vulnerabilities
The initial approach has vulnerabilities, particularly when malicious processes sniff the messages. An attacker could exploit knowledge of leader messages to manipulate the born list's size and attempt to gain control of the leadership.

# Objectives
The main objectives of this project are to enhance the security of the protocol using two techniques:

- Homomorphic Hashes: By replacing the size and identifier in leader messages with a homomorphic hash, we can verify if the receiver has an equal or larger born list. If it's a subset, no leadership change occurs. In case of equality, tie-breaking can be done using the identifier.

- Bloom Filters: Sending a bloom filter of the born list in addition to the homomorphic hash allows the receiver to check if its identifier is included in the sender's bloom filter. If it is, the receiver acknowledges the sender as the leader.

- ZK-SNARKs: Sending a ZK-proof for guarantee the veracity of the received message.

# Current Stage

Currently the project is under construction. It is implemented the original version and some communications and criptographic techniques, but it requires to review the state of the art about homomorphic hashes and bloom filter solid github projects. Regarding ZK-SNARK it is not designed neither introduced yet.

Looking for collaborators. All feedback appreciated.


# License
This project is licensed under the MIT License.

# Contact
For any questions or concerns, feel free to reach out to c.gomez.calzado@gmail.com.