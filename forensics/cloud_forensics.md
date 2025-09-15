Cloud Forensics

Overview:

This project showcases end-to-end forensic techniques in a simulated IaaS (Infrastructure-as-a-Service) environment. By combining cloud-native logging, threat detection, and evidence reconstruction, it walks through how analysts can investigate incidents in AWS. The approach blends red and blue team perspectives, from initial access to detection and remediation

Project Structure:

1. Cloud Environment Setup
Platform: AWS (Free Tier compatible) Purpose: Create a secure but realistic environment to simulate attacks 		and perform forensics.

Architecture
	•	1x public EC2 instance (Ubuntu) — vulnerable entry point
	•	1x private EC2 instance (internal server)
	•	S3 bucket with test data
	•	IAM roles with scoped permissions
	•	Logging Services:
	◦	AWS CloudTrail (API activity)
	◦	GuardDuty (threat detection)
	◦	VPC Flow Logs (network metadata)
	◦	S3 Server Access Logs



2. Incident Simulation:

Simulated Threat Scenario
	•	Attacker gains access via exposed SSH key (public EC2)
	•	Establishes reverse shell connection
	•	Uploads enumeration script (e.g. custom recon.sh)
	•	Exfiltrates S3 objects via stolen IAM permissions
	•	Performs lateral movement into private subnet (optional)



	Log Acquisition for Forensics
	•	Pulled logs from:
	◦	CloudTrail (API call history)
	◦	VPC Flow Logs (network metadata)
	◦	GuardDuty (threat intelligence and alerts)
	◦	S3 Access Logs (object access attempts)
	◦	Instance-level syslogs (via CloudWatch or agent forwarding)
Logs were exported and preserved in JSON/CSV for local analysis using command-line tools (jq, grep)


4. Timeline Reconstruction & Analysis
Rebuilding the attacker timeline using CloudTrail and VPC logs.
Key Findings:
	•	ConsoleLogin event from an unexpected region/IP
	•	DescribeInstances → enumeration of EC2 assets
	•	GetObject calls to S3 bucket following role compromise
	•	VPC logs showed reverse shell to attacker-controlled IP
# Example CloudTrail filter
jq '.Records[] | select(.eventName == "GetObject")' cloudtrail.json

5. Detection & Response Strategy
Response Simulation:
	•	Identified IAM role escalation via misconfigured trust policy
	•	Used GuardDuty alert to trace abuse pattern
	•	Isolated EC2 instance via security group change
	•	Rotated exposed IAM credentials
	•	Applied SCP (Service Control Policy) to restrict S3 access globally


Tools & Techniques Used

AWS CloudTrail
API activity logging (crucial for timeline)
AWS GuardDuty
Threat detection and anomaly spotting
AWS Config
Change tracking for compliance & rollback
SIFT Workstation
Cloud artifact parsing and investigation
Open-source Tools
cloud-forensics-utils, awslogs, jq

Red vs Blue Benefits

Test visibility gaps in cloud telemetry
Validate incident visibility end-to-end
Evade monitoring using indirect API chains
Detect misuse through behavior profiling (IAM, S3)
Abuse misconfigured roles or trust relationships
Harden least-privilege roles and rotation policies


Key Takeaways

This cloud forensics lab reinforced that visibility is everything in IaaS environments. Without logging (especially CloudTrail), forensics becomes guesswork. The simulation highlighted how forensic timelines help connect alerts to attacker behaviors.

