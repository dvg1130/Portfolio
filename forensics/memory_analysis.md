Memory Forensic Analysis with Volatility 3

Overview: This project demonstrates a full memory forensics pipeline leveraging a golden image as a baseline for analysis. Using live memory acquisition tools and Volatility 3, it showcases how defenders and adversaries alike can gain critical insight into system behavior captured in RAM. The workflow includes capturing memory dumps, building custom volatility profiles, and identifying indicators of compromise (IOCs) through comparative analysis.

Golden Image — Establishing the Baseline

	A Golden Image is an ideal, clean snapshot of a system configured to a known good state. It serves as a reliable baseline for forensic analysis and can be used as a cold or hot backup to quickly restore systems to a consistent, secure, and operational state.

Configuration Process

	•	Deployed a fresh Linux VM (Ubuntu 22.04) with only essential packages installed.
	•	Performed system hardening by disabling unnecessary services and applying security patches.
	•	Documented system state: kernel version, installed software, and configuration files.

Use Case Benefits

	•	Blue Team: Offers a trustworthy benchmark for detecting unauthorized changes in system memory.
	•	Red Team: Helps simulate clean environments for stealth testing and tool validation.



Memory Acquisition with LiME

Once the golden image was validated, memory was captured using LiME—a trusted kernel module for Linux memory acquisition.


Golden and Compromised Dumps

	•	Captured memory from the golden image to represent normal system behavior.
	•	Created a second dump from a compromised VM, preloaded with anomalies (e.g., backdoor listener, suspicious script).



Acquisition Options

```bash
# Local dump to file:

sudo insmod lime.ko "path=/home/user/memdump.lime format=lime"


# Remote dump via TCP:

# On remote receiver:
nc -l -p 4444 > memory_dump.lime

# On sending system:
sudo insmod lime.ko "path=tcp:<RECEIVER_IP>:4444 format=lime"


# LiME’s lime format includes metadata to improve parsing accuracy with Volatility.
```
Volatility 3 Analysis

Volatility 3 enables OS-agnostic forensic analysis by parsing memory images and extracting kernel-level artifacts using custom symbol profiles.

Forensic Plugins Used

```bash
vol -f memory.lime --symbol-file golden_profile.json linux.pslist     # Active processes

vol -f memory.lime --symbol-file golden_profile.json linux.malfind    # Suspicious memory regions

vol -f memory.lime --symbol-file golden_profile.json linux.lsmod      # Kernel modules

vol -f memory.lime --symbol-file golden_profile.json linux.netstat    # Network activity

vol -f memory.lime --symbol-file golden_profile.json linux.lsof       # Open files

vol -f memory.lime --symbol-file golden_profile.json linux.bash       # Bash history

vol -f memory.lime --symbol-file golden_profile.json linux.environ    # Environment variables
```

Comparative Analysis

To isolate anomalies, I ran identical plugins against both dumps and used command-line diff tools to flag deviations.

Example: Process Diff
```bash
vol -f golden.lime --symbol-file profile.json linux.pslist > golden_ps.txt
vol -f compromised.lime --symbol-file profile.json linux.pslist > compromised_ps.txt
diff golden_ps.txt compromised_ps.txt
```


Findings
	•	Detected hidden process injection not present in the golden image.
	•	Identified unauthorized reverse shell established via  over TCP.
	•	Highlighted a non-standard kernel module loaded at runtime.


Red & Blue Perspectives

	•	Blue Team - Validated memory forensics as a tool for post-compromise triage and IOC hunting

	•	Red Team - Observed detection surface left by common persistence methods and runtime loaders