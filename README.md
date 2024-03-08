ZK-leader-election
===
# Introduction
This project implements a distributed leader election algorithm based on the proposal by [Carlos Gómez-Calzado et al.](https://addi.ehu.es/bitstream/handle/10810/16424/TESIS_CARLOS_GOMEZ_CALZADO.pdf?sequence=1). The algorithm utilizes two types of messages: a "birth" message, broadcasted using the reliable broadcast primitive when a process starts, and a leader election criterion based on the size of the "born" list and, in case of a tie, the process identifier.

# Vulnerabilities
The initial approach has vulnerabilities, particularly when malicious processes sniff the messages. An attacker could exploit knowledge of leader messages to manipulate the born list's size and attempt to gain control of the leadership.

# Objectives
The main objectives of this project are to enhance the security of the protocol using two techniques:

1. Homomorphic Hashes:

- Objective Question: Can we enhance the security of the protocol by replacing the size and identifier in leader messages with a homomorphic hash, allowing us to verify if the receiver has an equal or larger born list? And in the case of a subset, can we ensure that no leadership change occurs, with tie-breaking based on the identifier in the event of equality?

2. ZK-SNARKs:

- Objective Question: Can we send a ZK-proof to guarantee the veracity of the received message as part of the project's objectives?

3. Transition to State-Machine Replication:

- Objective Question: Is it feasible to transition from a standard one-shot consensus to a paradigm of state-machine replication as part of the project's objectives?

# Current Stage

Currently, the project is under construction. The original version has been implemented along with some communication and cryptographic techniques. However, it requires a review of the state-of-the-art regarding homomorphic hashes and other solid GitHub projects related to cryptographic techniques. As for ZK-SNARK, it has not been designed or introduced yet.

We are actively seeking collaborators, and all feedback is appreciated.


# License
This project is licensed under the MIT License.

# Contact
For any questions or concerns, feel free to reach out to c.gomez.calzado@gmail.com.