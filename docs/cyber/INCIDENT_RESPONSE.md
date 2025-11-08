# Incident Response Playbook

**Version:** 1.0
**Last Updated:** 2025-11-05
**Classification:** CONFIDENTIAL - Internal Security Use Only

---

## Emergency Contacts

| Role | Contact | Phone | Email | Availability |
|------|---------|-------|-------|--------------|
| **Security Lead** | [Name] | [Phone] | security@univyn.com | 24/7 |
| **CTO** | [Name] | [Phone] | cto@univyn.com | 24/7 |
| **DevOps Lead** | [Name] | [Phone] | devops@univyn.com | 24/7 |
| **Legal** | [Name] | [Phone] | legal@univyn.com | Business hours |
| **PR/Communications** | [Name] | [Phone] | pr@univyn.com | Business hours |

---

## Incident Severity Levels

| Level | Description | Response Time | Examples |
|-------|-------------|---------------|----------|
| **P0 - Critical** | Active breach, data exposure, complete service outage | Immediate (< 15 min) | Active data breach, RCE exploit, database compromised |
| **P1 - High** | Potential breach, significant security risk, major functionality impaired | < 1 hour | Credentials leaked in logs, SQLi vulnerability discovered |
| **P2 - Medium** | Security incident with contained impact | < 4 hours | Brute force attack detected, suspicious account activity |
| **P3 - Low** | Minor security issue, no immediate risk | < 24 hours | Outdated dependency with low-risk CVE |

---

## General Incident Response Process

```
[Detection] → [Triage] → [Containment] → [Eradication] → [Recovery] → [Post-Incident Review]
     ↓            ↓            ↓              ↓               ↓               ↓
  Monitoring   Assess    Block/Isolate   Fix Root      Restore        Document
  Alerts      Severity   Threat          Cause         Service        Lessons
```

### Phase 1: Detection (0-5 minutes)

**Sources:**
- Automated monitoring alerts
- Security scan findings
- User reports
- Log analysis
- External notification (security researcher, customer)

**Actions:**
1. Document initial report with timestamp
2. Gather preliminary information
3. Create incident ticket in tracking system
4. Alert on-call security engineer

### Phase 2: Triage (5-15 minutes)

**Actions:**
1. Verify incident is real (not false positive)
2. Assess severity level (P0-P3)
3. Identify affected systems/data
4. Determine if incident is ongoing
5. Assemble incident response team based on severity
6. Establish communication channel (Slack #incident-response or dedicated channel)

**Severity Assessment Questions:**
- Is data being actively exfiltrated?
- Are credentials compromised?
- Is service availability impacted?
- Are multiple users/systems affected?
- Is the vulnerability being actively exploited?

### Phase 3: Containment (15 minutes - 2 hours)

**Immediate Actions:**
- Prevent further damage
- Preserve evidence
- Maintain service availability if possible

**Proceed to specific playbook based on incident type** (see below)

### Phase 4: Eradication (varies)

**Actions:**
1. Identify root cause
2. Remove threat/vulnerability
3. Patch systems
4. Rotate compromised credentials
5. Update security controls
6. Verify threat is eliminated

### Phase 5: Recovery (varies)

**Actions:**
1. Restore systems from clean backups if needed
2. Verify system integrity
3. Monitor for reoccurrence
4. Gradually restore service
5. Communicate with stakeholders

### Phase 6: Post-Incident Review (within 1 week)

**Actions:**
1. Document complete timeline
2. Identify lessons learned
3. Update security controls
4. Update incident response procedures
5. Conduct team retrospective
6. Create action items with owners

---

## Playbook 1: Credentials Compromised

### Scenario
Database credentials, API keys, or secrets exposed in:
- Git repository
- Logs
- Error messages
- Third-party breach

### Immediate Actions (15-30 minutes)

1. **Rotate Compromised Secret IMMEDIATELY**
   ```bash
   # Generate new key
   NEW_KEY=$(openssl rand -base64 32)

   # Update in production secret manager
   # (AWS Secrets Manager example)
   aws secretsmanager update-secret \
       --secret-id connect/production/token-key \
       --secret-string "$NEW_KEY"

   # Restart application to pick up new secret
   kubectl rollout restart deployment/connect-api
   ```

2. **Revoke All Active Sessions** (if TOKEN_SYMMETRIC_KEY compromised)
   ```sql
   DELETE FROM sessions WHERE expires_at > NOW();
   ```

3. **Assess Exposure**
   - How was credential exposed?
   - For how long?
   - Who had access?
   - Any evidence of unauthorized use?

4. **Check for Unauthorized Access**
   ```sql
   -- Check for suspicious login activity
   SELECT * FROM login_attempts
   WHERE attempted_at > '[EXPOSURE_TIME]'
     AND attempt_result = 'success'
   ORDER BY attempted_at DESC;

   -- Check for unusual database activity (if DB creds compromised)
   -- Review database audit logs
   ```

5. **Block Suspicious IPs** (if unauthorized access detected)
   ```bash
   # Add to firewall rules
   # Example with iptables
   iptables -A INPUT -s [SUSPICIOUS_IP] -j DROP

   # Or update WAF/load balancer rules
   ```

### Follow-up Actions (1-24 hours)

6. **Remove Secret from Exposed Location**
   - If in Git: Use git-filter-repo to remove from history
   - If in logs: Purge logs, update log rotation
   - If in error messages: Fix code, deploy update

7. **Notify Stakeholders**
   - Security team
   - Management
   - Legal (if customer data potentially accessed)
   - Affected users (if required by regulation)

8. **Investigate Root Cause**
   - Why was secret exposed?
   - How can we prevent this in future?

### Recovery Actions (1-7 days)

9. **Implement Preventive Controls**
   - Add pre-commit hooks to prevent secret commits
   - Enable GitHub secret scanning
   - Implement secret scanning in CI/CD
   - Update logging to redact secrets

10. **Documentation**
    - Complete incident report
    - Update runbooks
    - Conduct post-incident review

### Evidence Preservation

- Git commit hashes (if in repo)
- Log snippets (redacted)
- Timeline of events
- Access logs during exposure window

---

## Playbook 2: Account Takeover / Unauthorized Access

### Scenario
- User reports account compromise
- Suspicious activity detected in logs
- Multiple accounts accessed from same unusual location/IP

### Immediate Actions (15-30 minutes)

1. **Lock Affected Account(s)**
   ```sql
   UPDATE users
   SET is_locked = TRUE, locked_until = NULL
   WHERE id IN ('[USER_IDS]');
   ```

2. **Revoke All Sessions for Affected User(s)**
   ```sql
   DELETE FROM sessions
   WHERE user_id IN ('[USER_IDS]');
   ```

3. **Analyze Attack Pattern**
   ```sql
   -- Get login history for affected user
   SELECT * FROM login_attempts
   WHERE user_id = '[USER_ID]'
     AND attempted_at > NOW() - INTERVAL '7 days'
   ORDER BY attempted_at DESC;

   -- Check for credential stuffing pattern
   SELECT ip_address, COUNT(*) as attempt_count
   FROM login_attempts
   WHERE attempt_result = 'failed_password'
     AND attempted_at > NOW() - INTERVAL '1 hour'
   GROUP BY ip_address
   HAVING COUNT(*) > 10;
   ```

4. **Block Malicious IPs** (if pattern identified)
   ```bash
   # Add to rate limiting or firewall
   # Update WAF rules to block IPs
   ```

5. **Assess Damage**
   - What actions did attacker perform?
   - What data was accessed?
   - Were any changes made?

### Follow-up Actions (1-4 hours)

6. **Notify Affected User(s)**
   - Email notification
   - Explain what happened
   - Provide unlock/password reset instructions
   - Security recommendations

7. **Force Password Reset**
   - Generate password reset token
   - Send to user's verified email
   - Require strong new password

8. **Review Account Activity**
   - Check for unauthorized changes
   - Review access to sensitive data
   - Check for malicious content posted

9. **Investigate Attack Vector**
   - Credential stuffing?
   - Phishing?
   - Session hijacking?
   - Database breach?

### Recovery Actions (1-7 days)

10. **Strengthen Controls**
    - If credential stuffing: Enhance rate limiting
    - If phishing: User education, add MFA
    - If session hijacking: Review session security

11. **Monitor for Reoccurrence**
    - Alert on suspicious activity from affected accounts
    - Monitor for pattern across other accounts

### Evidence Preservation

- Login attempt logs
- IP addresses involved
- User actions during compromise window
- Email communications with user

---

## Playbook 3: SQL Injection or Code Injection Attack

### Scenario
- SQL injection attempt detected in logs
- Unusual database queries detected
- Error messages indicate injection attempt

### Immediate Actions (IMMEDIATE - 15 minutes)

1. **Verify Attack is Real**
   - Check logs for injection patterns
   - Example SQLi patterns: `' OR '1'='1`, `'; DROP TABLE`, `UNION SELECT`
   - Example XSS: `<script>`, `javascript:`, `onerror=`

2. **Block Attacker IP Immediately**
   ```bash
   # Add to firewall
   iptables -A INPUT -s [ATTACKER_IP] -j DROP

   # Or use WAF rules
   ```

3. **Assess if Attack Succeeded**
   ```sql
   -- Check for unusual database modifications
   SELECT * FROM [audit_table]
   WHERE modified_at > '[ATTACK_TIME]';

   -- Check for data exfiltration attempts
   -- Review database logs for unusual queries
   ```

4. **If Attack Succeeded (Data Breach):**
   - **ESCALATE TO P0 IMMEDIATELY**
   - Notify Security Lead and CTO
   - Preserve all evidence
   - Consider taking system offline if active

### Follow-up Actions (1-4 hours)

5. **Identify Vulnerable Endpoint**
   - Review code for injection point
   - Check if parameterized queries used
   - Identify affected endpoints

6. **Patch Vulnerability**
   - Fix code to use parameterized queries
   - Add input validation
   - Deploy fix immediately

7. **Verify Fix**
   - Test injection attempts against patched code
   - Conduct code review
   - Run automated security tests

8. **Check for Other Vulnerable Endpoints**
   ```bash
   # Search for potential SQL injection vulnerabilities
   grep -r "fmt.Sprintf.*SELECT" .
   grep -r "string concatenation" .

   # Run SAST tools
   gosec ./...
   ```

### Recovery Actions (1-7 days)

9. **Comprehensive Code Audit**
   - Review all database interaction code
   - Ensure SQLC used consistently
   - Add security-focused code review process

10. **Enhanced Monitoring**
    - Add alerting for injection patterns in logs
    - Implement WAF rules to detect injection attempts
    - Regular security scanning

### Evidence Preservation

- Request logs showing injection attempts
- IP address(es) of attacker
- Vulnerable code (before fix)
- Database logs during attack window

---

## Playbook 4: Brute Force Attack

### Scenario
- High volume of failed login attempts
- Account lockout alerts triggered
- Rate limiting violations

### Immediate Actions (5-15 minutes)

1. **Verify Attack Pattern**
   ```sql
   -- Check failed login attempts
   SELECT ip_address, COUNT(*) as attempts
   FROM login_attempts
   WHERE attempt_result IN ('failed_password', 'failed_user_not_found')
     AND attempted_at > NOW() - INTERVAL '15 minutes'
   GROUP BY ip_address
   HAVING COUNT(*) > 10
   ORDER BY attempts DESC;
   ```

2. **Block Attacking IPs**
   ```bash
   # Get list of IPs with excessive failed attempts
   # Add to firewall or WAF blocklist

   # Example: Block top 10 offending IPs
   for ip in $(psql -t -c "SELECT DISTINCT ip_address FROM login_attempts WHERE..."); do
       iptables -A INPUT -s $ip -j DROP
   done
   ```

3. **Verify Account Lockout Working**
   ```sql
   -- Check locked accounts
   SELECT username, is_locked, locked_until, failed_login_attempts
   FROM users
   WHERE is_locked = TRUE;
   ```

4. **Check if Any Accounts Compromised**
   ```sql
   -- Look for successful logins from attacking IPs
   SELECT DISTINCT user_id, username, ip_address
   FROM login_attempts
   WHERE ip_address IN (SELECT ip_address FROM login_attempts WHERE ...)
      AND attempt_result = 'success';
   ```

### Follow-up Actions (1-4 hours)

5. **Analyze Attack Details**
   - Single IP or distributed?
   - Targeting specific accounts or spray attack?
   - Password patterns (if visible in logs)

6. **Notify Targeted Users** (if specific accounts targeted)
   - Email alert
   - Recommend password change if weak
   - Offer account security review

7. **Adjust Rate Limiting** (if attack bypassed limits)
   - Lower authentication rate limit threshold
   - Implement progressive delays
   - Add CAPTCHA for repeated failures

### Recovery Actions (1-7 days)

8. **Review and Enhance Controls**
   - Analyze effectiveness of account lockout
   - Consider adding CAPTCHA
   - Implement MFA for high-risk accounts
   - Add geo-blocking if attack from specific region

9. **Monitor for Continued Attacks**
   - Alert on similar patterns
   - Track attacking IP ranges
   - Coordinate with upstream providers if needed

### Evidence Preservation

- Login attempt logs (full data for attack window)
- IP addresses involved
- Usernames targeted
- Success rate of attacks

---

## Playbook 5: Data Breach / Unauthorized Data Access

### Scenario
- Evidence of unauthorized data access
- Data exfiltration detected
- External notification of data exposure

### Immediate Actions (IMMEDIATE - P0)

1. **Activate P0 Incident Response**
   - Notify Security Lead, CTO, CEO, Legal
   - Assemble crisis response team
   - Establish communication channel

2. **Contain the Breach**
   - If ongoing, consider taking affected system offline
   - Block attacker access immediately
   - Preserve all evidence

3. **Assess Scope**
   - What data was accessed/exfiltrated?
   - How many users affected?
   - What PII/sensitive data involved?
   - Time window of unauthorized access?

4. **Preserve Evidence**
   - Take system snapshots
   - Copy logs to secure location
   - Document all actions taken
   - Do NOT modify affected systems unless necessary

### Follow-up Actions (1-4 hours)

5. **Identify Attack Vector**
   - How did attacker gain access?
   - What vulnerability was exploited?
   - Were credentials compromised?

6. **Patch Vulnerability**
   - Fix immediate security issue
   - Deploy patch
   - Verify patch effectiveness

7. **Notify Stakeholders**
   - Management
   - Legal team
   - Prepare for potential regulatory notification

### Recovery Actions (1-72 hours)

8. **Regulatory Notification** (if required)
   - **GDPR**: 72 hours to notify supervisory authority
   - **Other regulations**: Check requirements
   - Legal team leads this process

9. **User Notification** (if PII exposed)
   - Draft notification with Legal/PR
   - Email affected users
   - Provide recommended actions
   - Offer support resources

10. **Forensic Analysis**
    - Detailed investigation
    - Document complete attack chain
    - Identify all affected data
    - External forensics if needed

11. **Remediation**
    - Fix all identified vulnerabilities
    - Enhance security controls
    - Conduct security audit
    - Implement additional monitoring

### Post-Incident Actions (1-4 weeks)

12. **Comprehensive Security Review**
    - Third-party security assessment
    - Penetration testing
    - Security architecture review

13. **Public Communication** (if needed)
    - Press release (coordinated with PR)
    - FAQ for customers
    - Transparency report

14. **Regulatory Compliance**
    - Respond to authority inquiries
    - Provide evidence of remediation
    - Update compliance documentation

### Evidence Preservation (CRITICAL)

- Complete system snapshots
- All logs (do not delete)
- Database dumps (if safe)
- Network traffic captures
- Timeline of all events
- Chain of custody documentation

---

## Communication Templates

### Internal Alert (P0/P1)

```
🚨 SECURITY INCIDENT - P[LEVEL]

Incident ID: INC-YYYYMMDD-NNN
Severity: P[LEVEL]
Status: [ACTIVE/CONTAINED/RESOLVED]
Incident Lead: [Name]

Summary: [Brief description]

Affected Systems: [List]
Impact: [Description]
Current Actions: [What's being done]

Next Update: [Time]

Response Team: [List team members]
Communication Channel: #incident-[ID]
```

### User Notification (Account Compromise)

```
Subject: Security Alert: Suspicious Activity on Your Account

Dear [User],

We detected suspicious login attempts on your Univyn account on [Date].
As a precautionary measure, we have temporarily locked your account.

What happened:
[Brief, non-technical explanation]

What we did:
- Locked your account to prevent unauthorized access
- Terminated all active sessions
- Investigated the incident

What you should do:
1. Reset your password using this link: [Link]
2. Review your recent account activity
3. Enable two-factor authentication (recommended)
4. Use a unique, strong password

If you did not attempt to log in from [Location] on [Date], please contact
our security team immediately at security@univyn.com.

We take your security seriously and apologize for any inconvenience.

Sincerely,
Univyn Security Team
```

### Breach Notification (GDPR Compliant)

```
Subject: Important Security Notice Regarding Your Univyn Account

Dear [User],

We are writing to inform you of a security incident that may have affected
your personal information.

What happened:
[Detailed explanation of breach]

What information was involved:
[List specific data types: name, email, etc.]

What we have done:
[List remediation steps]

What you should do:
[Specific recommended actions]

Additional information:
[Support resources, FAQs, contact information]

We sincerely apologize for this incident and any concern it may cause.
The security of your information is our top priority.

For questions, please contact: security@univyn.com

Sincerely,
[Name]
[Title]
Univyn
```

---

## Post-Incident Review Template

**Incident ID:** INC-YYYYMMDD-NNN
**Incident Type:** [Credential Compromise/Breach/Attack/etc.]
**Severity:** P[0-3]
**Date:** YYYY-MM-DD
**Incident Lead:** [Name]

### Timeline

| Time | Event | Action Taken | Owner |
|------|-------|--------------|-------|
| HH:MM | Incident detected | Alert received | |
| HH:MM | Triage completed | Severity assessed as P1 | |
| HH:MM | Containment | [Action] | |
| HH:MM | Resolution | [Action] | |

### Impact Assessment

- **Systems Affected:** [List]
- **Data Affected:** [Description]
- **Users Affected:** [Count/description]
- **Service Availability:** [Uptime %]
- **Duration:** [Total time from detection to resolution]

### Root Cause

[Detailed explanation of what caused the incident]

### What Went Well

- [List things that worked well during response]

### What Could Be Improved

- [List areas for improvement]

### Action Items

| Action | Owner | Priority | Due Date | Status |
|--------|-------|----------|----------|--------|
| [Action 1] | [Name] | High | YYYY-MM-DD | Open |
| [Action 2] | [Name] | Medium | YYYY-MM-DD | Open |

### Lessons Learned

- [Key lessons]

### Updated Procedures

- [List any procedure updates resulting from this incident]

---

## Incident Log

**All incidents must be logged in the incident tracking system.**

Required information:
- Incident ID
- Date/time detected
- Severity level
- Type of incident
- Systems affected
- Response team
- Actions taken
- Resolution
- Post-incident review completion

---

## Regular Drills

**Frequency:** Quarterly

**Types:**
1. **Tabletop Exercise**: Team discusses hypothetical incident (1 hour)
2. **Simulated Incident**: Technical drill with test system (2-4 hours)
3. **Full Response Exercise**: Complete incident simulation with all teams (Half day)

**Next Scheduled Drill:** [Date]

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-05 | Security Team | Initial incident response playbook |

**Next Review Date:** 2026-02-05 (Quarterly)

**Distribution:** Security Team, Engineering Leadership, DevOps

---

**For active incidents, contact:**
**security@univyn.com** or **[On-Call Phone]**

**This document contains sensitive security information. Distribute only to authorized personnel.**
