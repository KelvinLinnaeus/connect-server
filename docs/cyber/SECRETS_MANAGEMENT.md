# Secrets Management Guide

**Version:** 1.0
**Last Updated:** 2025-11-05
**Classification:** Internal Use

## Overview

This document describes how to securely manage secrets, credentials, and sensitive configuration for the Univyn Connect Server. Proper secrets management is critical to prevent unauthorized access and data breaches.

---

## Table of Contents

1. [Critical Secrets](#critical-secrets)
2. [Environment Setup](#environment-setup)
3. [Secret Generation](#secret-generation)
4. [Storage & Distribution](#storage--distribution)
5. [Secret Rotation](#secret-rotation)
6. [Development vs Production](#development-vs-production)
7. [Security Best Practices](#security-best-practices)
8. [Incident Response](#incident-response)

---

## Critical Secrets

### Required Secrets

| Secret | Purpose | Format | Sensitivity |
|--------|---------|--------|-------------|
| `DATABASE_URL` | PostgreSQL connection | Connection string | **CRITICAL** |
| `TOKEN_SYMMETRIC_KEY` | PASETO token encryption | ≥32 char string | **CRITICAL** |
| `POSTGRES_PASSWORD` | Local dev DB password | String | **HIGH** |
| `TEST_DATABASE_URL` | Test DB connection | Connection string | **MEDIUM** |

### Optional Secrets

| Secret | Purpose | Sensitivity |
|--------|---------|-------------|
| Third-party API keys | External integrations | **HIGH** |
| Email service credentials | Notification delivery | **MEDIUM** |
| Cloud storage credentials | File uploads | **HIGH** |

---

## Environment Setup

### 1. Copy Environment Template

```bash
# Copy the example environment file
cp .env.example .env

# OR use app.env (both are supported)
cp .env.example app.env
```

### 2. Generate Required Secrets

See [Secret Generation](#secret-generation) section below.

### 3. Configure Environment Variables

Edit `.env` or `app.env` and replace all `CHANGE_ME` placeholders:

```bash
# Edit with your preferred editor
nano .env
# or
vim .env
```

### 4. Verify Configuration

Ensure all required variables are set:

```bash
# Check that required variables are present
grep -E "DATABASE_URL|TOKEN_SYMMETRIC_KEY" .env
```

### 5. Set File Permissions (Linux/Mac)

Protect your environment files from unauthorized access:

```bash
chmod 600 .env
chmod 600 app.env
```

---

## Secret Generation

### TOKEN_SYMMETRIC_KEY

**Requirement:** Must be at least 32 characters long (enforced by application)

**Generate using OpenSSL (Linux/Mac):**
```bash
openssl rand -base64 32
```

**Generate using PowerShell (Windows):**
```powershell
$bytes = New-Object byte[] 32
[Security.Cryptography.RNGCryptoServiceProvider]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes)
```

**Generate using Go:**
```bash
go run -c 'package main
import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)
func main() {
    b := make([]byte, 32)
    rand.Read(b)
    fmt.Println(base64.StdEncoding.EncodeToString(b))
}'
```

**Example Output:**
```
3fX9kL2mN8pQ5rS7tU0vW1xY3zA6bC4dE5fG7hI9jK2lM4nO6pQ8rS0tU2vW4xY6z
```

### DATABASE_URL Password Component

**Requirement:** Strong, random password (minimum 16 characters)

**Generate Strong Password (Linux/Mac):**
```bash
# 32-character alphanumeric + special chars
openssl rand -base64 24 | tr -d "=+/" | cut -c1-32
```

**Generate Strong Password (PowerShell):**
```powershell
Add-Type -AssemblyName System.Web
[System.Web.Security.Membership]::GeneratePassword(32,8)
```

**Password Requirements:**
- Minimum 16 characters (32+ recommended)
- Mix of uppercase, lowercase, numbers, special characters
- Avoid dictionary words
- Unique per environment

### POSTGRES_PASSWORD

Use same generation method as DATABASE_URL password above.

---

## Storage & Distribution

### Development Environment

**Option 1: Local .env File (Recommended for Individual Developers)**

```bash
# Create .env file (gitignored)
cp .env.example .env

# Edit with your credentials
nano .env
```

**Option 2: Shell Environment Variables**

```bash
# Add to ~/.bashrc or ~/.zshrc
export DATABASE_URL="postgresql://user:pass@localhost:5432/connect?sslmode=disable"
export TOKEN_SYMMETRIC_KEY="your-32-char-or-longer-key"

# Reload shell
source ~/.bashrc
```

### Production Environment

**❌ NEVER:**
- Commit secrets to Git
- Store secrets in code
- Share secrets via email/chat
- Store secrets in plain text files without encryption
- Use the same secrets across environments

**✅ ALWAYS:**
- Use a secrets management service
- Encrypt secrets at rest
- Audit secret access
- Rotate secrets regularly
- Use unique secrets per environment

**Recommended Secret Management Solutions:**

1. **Cloud Provider Secret Managers**
   - AWS Secrets Manager
   - Google Cloud Secret Manager
   - Azure Key Vault

2. **Self-Hosted Solutions**
   - HashiCorp Vault
   - Docker Secrets (for Docker Swarm)
   - Kubernetes Secrets

3. **CI/CD Integration**
   - GitHub Secrets (for GitHub Actions)
   - GitLab CI/CD Variables
   - CircleCI Environment Variables

### Example: AWS Secrets Manager

```bash
# Store secret
aws secretsmanager create-secret \
    --name connect/production/database-url \
    --secret-string "postgresql://user:pass@host:5432/db?sslmode=require"

# Retrieve secret in application startup script
export DATABASE_URL=$(aws secretsmanager get-secret-value \
    --secret-id connect/production/database-url \
    --query SecretString \
    --output text)
```

### Example: Docker Secrets

```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    image: connect-api:latest
    secrets:
      - database_url
      - token_key
    environment:
      DATABASE_URL_FILE: /run/secrets/database_url
      TOKEN_SYMMETRIC_KEY_FILE: /run/secrets/token_key

secrets:
  database_url:
    external: true
  token_key:
    external: true
```

---

## Secret Rotation

### Why Rotate Secrets?

- Limit exposure window if compromised
- Comply with security policies
- Reduce impact of insider threats
- Meet compliance requirements (SOC 2, PCI DSS)

### Rotation Schedule

| Secret | Recommended Frequency | Trigger Events |
|--------|-----------------------|----------------|
| `DATABASE_URL` | Every 90 days | Suspected compromise, employee departure |
| `TOKEN_SYMMETRIC_KEY` | Every 90 days | Suspected compromise, major security incident |
| `POSTGRES_PASSWORD` | Every 90 days | Developer laptop lost/stolen |

### TOKEN_SYMMETRIC_KEY Rotation Procedure

**⚠️ WARNING:** Rotating this key will invalidate all existing tokens. Plan for coordinated deployment.

**Step 1: Generate New Key**
```bash
NEW_KEY=$(openssl rand -base64 32)
echo "New key: $NEW_KEY"
```

**Step 2: Update Code for Dual-Key Support**

```go
// internal/util/config.go
type Config struct {
    // ... existing fields
    TokenSymmetricKey    string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
    TokenSymmetricKeyOld string `mapstructure:"TOKEN_SYMMETRIC_KEY_OLD"` // Add this
}

// internal/service/auth/token_maker.go
// Update VerifyToken to try both keys
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
    payload, err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload)
    if err != nil {
        // Try old key if verification fails
        if maker.symmetricKeyOld != "" {
            payload, err = maker.paseto.Decrypt(token, maker.symmetricKeyOld, &payload)
        }
    }
    // ... rest of verification
}
```

**Step 3: Deploy with Both Keys**

```bash
# Update environment with both keys
export TOKEN_SYMMETRIC_KEY="<new-key>"
export TOKEN_SYMMETRIC_KEY_OLD="<old-key>"

# Deploy application
make deploy
```

**Step 4: Wait for Token Expiry**

Wait for maximum refresh token duration (24 hours by default). All old tokens will naturally expire.

**Step 5: Remove Old Key**

```bash
# Remove old key from environment
unset TOKEN_SYMMETRIC_KEY_OLD

# Redeploy
make deploy
```

**Step 6: Update Code to Remove Dual-Key Support**

Remove `TOKEN_SYMMETRIC_KEY_OLD` from code and configuration.

### DATABASE_URL Rotation Procedure

**Step 1: Create New Database User**

```sql
-- Connect to database as admin
psql -h your-db-host -U admin -d connect

-- Create new user with strong password
CREATE USER connect_user_new WITH PASSWORD 'new-strong-password';

-- Grant same permissions as old user
GRANT ALL PRIVILEGES ON DATABASE connect TO connect_user_new;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO connect_user_new;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO connect_user_new;
```

**Step 2: Update DATABASE_URL**

```bash
# Update environment variable
export DATABASE_URL="postgresql://connect_user_new:new-strong-password@host:5432/connect?sslmode=require"
```

**Step 3: Deploy with Zero-Downtime Strategy**

```bash
# Rolling update (no downtime)
make deploy-rolling
```

**Step 4: Verify Connectivity**

```bash
# Check application logs
tail -f /var/log/connect/app.log | grep "database"

# Test database connection
make test-db-connection
```

**Step 5: Revoke Old Credentials**

```sql
-- After verifying new user works, revoke old user
DROP USER connect_user_old;
```

### Emergency Rotation (Suspected Compromise)

**IMMEDIATE ACTIONS:**

1. **Rotate Compromised Secret Immediately**
   ```bash
   # Generate new key
   NEW_KEY=$(openssl rand -base64 32)

   # Update production environment (use your secret manager)
   aws secretsmanager update-secret \
       --secret-id connect/production/token-key \
       --secret-string "$NEW_KEY"

   # Force application restart
   kubectl rollout restart deployment/connect-api
   ```

2. **Revoke All Active Sessions**
   ```sql
   -- Revoke all refresh tokens
   DELETE FROM sessions WHERE expires_at > NOW();
   ```

3. **Notify Security Team & Users**
   - Alert security team immediately
   - Log incident details
   - Notify affected users if necessary
   - Document timeline and actions taken

4. **Investigate & Remediate**
   - Review access logs
   - Identify compromise vector
   - Patch vulnerabilities
   - Update security controls

---

## Development vs Production

### Development Environment

**Characteristics:**
- Local database (Docker or local PostgreSQL)
- Relaxed CORS policy
- SSL/TLS optional
- Verbose logging
- Simple secrets (but still not hardcoded!)

**Example .env:**
```bash
ENVIRONMENT=development
SERVER_ADDRESS=0.0.0.0:8080
DATABASE_URL=postgresql://postgres:dev_password@localhost:5432/connect?sslmode=disable
TOKEN_SYMMETRIC_KEY=dev-key-at-least-32-characters-long-random
CORS_ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT=100
RATE_LIMIT_AUTH=10
```

### Staging Environment

**Characteristics:**
- Cloud-hosted database
- Production-like configuration
- SSL/TLS required
- Moderate logging
- Strong secrets

**Example Configuration:**
```bash
ENVIRONMENT=staging
SERVER_ADDRESS=0.0.0.0:8080
DATABASE_URL=postgresql://connect_staging:strong_random_pass@staging-db.internal:5432/connect?sslmode=require
TOKEN_SYMMETRIC_KEY=<32+ char random key different from prod>
CORS_ALLOWED_ORIGINS=https://staging.univyn.com
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT=100
RATE_LIMIT_AUTH=5
```

### Production Environment

**Characteristics:**
- Cloud-hosted database with encryption
- Strict CORS policy
- SSL/TLS enforced
- Structured logging with PII redaction
- Strong, unique secrets

**Example Configuration:**
```bash
ENVIRONMENT=production
SERVER_ADDRESS=0.0.0.0:8080
DATABASE_URL=postgresql://connect_prod:very_strong_random_32plus_char_pass@prod-db.internal:5432/connect?sslmode=verify-full
TOKEN_SYMMETRIC_KEY=<32+ char random key unique to production>
CORS_ALLOWED_ORIGINS=https://app.univyn.com,https://univyn.com
RATE_LIMIT_ENABLED=true
RATE_LIMIT_DEFAULT=100
RATE_LIMIT_AUTH=5
```

---

## Security Best Practices

### ✅ DO

1. **Use Strong, Random Secrets**
   - Generate with cryptographically secure methods
   - Minimum 32 characters for symmetric keys
   - Minimum 16 characters for passwords

2. **Unique Secrets Per Environment**
   - Development ≠ Staging ≠ Production
   - Each developer can have unique local credentials

3. **Encrypt Secrets at Rest**
   - Use secret management services
   - Encrypt backup files containing secrets
   - Use encrypted channels for distribution

4. **Audit Secret Access**
   - Log who accessed which secrets when
   - Review audit logs regularly
   - Alert on suspicious access patterns

5. **Rotate Regularly**
   - Follow rotation schedule
   - Document rotation procedures
   - Test rotation in staging first

6. **Least Privilege**
   - Database users should have minimum required permissions
   - Application should not run as root/admin
   - Limit secret access to necessary personnel only

7. **Version Control Protection**
   - Verify .env and app.env are in .gitignore
   - Use pre-commit hooks to prevent secret commits
   - Scan commits for accidentally committed secrets

### ❌ DON'T

1. **Never Commit Secrets to Git**
   - Check .gitignore includes .env, app.env
   - Use git-secrets or similar tools
   - If accidentally committed, rotate immediately

2. **Never Share Secrets via Insecure Channels**
   - No email, Slack, SMS, or unencrypted chat
   - Use dedicated secret sharing tools (1Password, LastPass)
   - Prefer secret management services

3. **Never Log Secrets**
   - Redact passwords from database connection strings
   - Mask tokens in logs
   - Never log request headers containing Authorization

4. **Never Use Default or Weak Secrets**
   - No "password", "123456", "admin"
   - No short or predictable values
   - No secrets derived from public information

5. **Never Reuse Secrets**
   - Across environments (dev/staging/prod)
   - Across applications
   - After rotation (don't rotate back to old value)

---

## Incident Response

### If Secrets Are Compromised

**SEVERITY: CRITICAL**

**Immediate Actions (within 1 hour):**

1. **Rotate Compromised Secret**
   - Follow emergency rotation procedure
   - Use new, strong random value

2. **Revoke Active Sessions**
   - Invalidate all user sessions if auth secrets compromised
   - Force re-authentication

3. **Notify Security Team**
   - Alert security@univyn.com (or your security contact)
   - Document incident details

4. **Block Unauthorized Access**
   - Check for suspicious activity in logs
   - Block malicious IPs if identified
   - Increase monitoring and alerting

**Short-term Actions (within 24 hours):**

5. **Investigate Compromise Vector**
   - How were secrets exposed?
   - Git history, logs, developer machines?
   - External breach or insider threat?

6. **Assess Impact**
   - Which secrets were exposed?
   - For how long were they exposed?
   - What data could have been accessed?

7. **Notify Stakeholders**
   - Management
   - Legal/compliance team
   - Affected users (if necessary per GDPR/regulations)

**Long-term Actions (within 1 week):**

8. **Remediate Root Cause**
   - Fix vulnerability that led to exposure
   - Implement additional security controls
   - Update security policies and training

9. **Post-Incident Review**
   - Document timeline and actions
   - Identify lessons learned
   - Update incident response procedures

10. **Compliance Reporting**
    - GDPR breach notification (if applicable, within 72 hours)
    - Other regulatory requirements

### If Secrets Are Accidentally Committed to Git

**SEVERITY: HIGH**

**Immediate Actions:**

1. **Rotate Secret Immediately**
   - Do not wait to clean Git history
   - Assume secret is compromised

2. **Remove from Git History**
   ```bash
   # Use git-filter-repo (preferred) or BFG Repo-Cleaner
   git filter-repo --path .env --invert-paths

   # Force push (coordinate with team)
   git push --force --all
   ```

3. **Notify Team**
   - All developers must re-clone repository
   - Clear local caches

4. **Verify Removal**
   - Check all branches and tags
   - Verify secret no longer appears in `git log --all -S "secret_value"`

**Prevention:**

- Use pre-commit hooks (see below)
- Enable GitHub secret scanning
- Regular security training

### Pre-commit Hook to Prevent Secret Commits

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
# Prevent committing files with potential secrets

FILES_PATTERN='\.(env|env\..*|pem|key|p12)$'
FORBIDDEN='(password|secret|key|token|credential|DATABASE_URL).*=.+'

# Check staged files
if git diff --cached --name-only | grep -E "$FILES_PATTERN"; then
    echo "ERROR: Attempting to commit sensitive files (.env, .key, etc.)"
    echo "These files should be in .gitignore"
    exit 1
fi

# Check content for secrets
if git diff --cached | grep -iE "$FORBIDDEN"; then
    echo "ERROR: Potential secret detected in commit"
    echo "Please remove secrets and use environment variables"
    exit 1
fi

exit 0
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Verification Checklist

Before deploying to any environment, verify:

- [ ] All `CHANGE_ME` placeholders replaced
- [ ] `TOKEN_SYMMETRIC_KEY` is at least 32 random characters
- [ ] `DATABASE_URL` uses strong password (16+ characters)
- [ ] `DATABASE_URL` uses SSL/TLS (sslmode=require) in production
- [ ] `CORS_ALLOWED_ORIGINS` is set to specific domains (NOT "*") in production
- [ ] All secrets are unique per environment
- [ ] Secrets are not committed to version control
- [ ] `.env` and `app.env` are in `.gitignore`
- [ ] File permissions are restrictive (600 on Linux/Mac)
- [ ] Secret access is logged and audited
- [ ] Rotation schedule is documented and automated
- [ ] Incident response procedure is documented and tested

---

## Additional Resources

### Tools

- **git-secrets**: Prevents committing secrets to Git
  - https://github.com/awslabs/git-secrets

- **truffleHog**: Finds secrets in Git history
  - https://github.com/trufflesecurity/trufflehog

- **detect-secrets**: Detects secrets in code
  - https://github.com/Yelp/detect-secrets

### Documentation

- **OWASP Secrets Management Cheat Sheet**
  - https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_CheatSheet.html

- **NIST Special Publication 800-57: Key Management**
  - https://csrc.nist.gov/publications/detail/sp/800-57-part-1/rev-5/final

### Secret Management Services

- **AWS Secrets Manager**: https://aws.amazon.com/secrets-manager/
- **Google Cloud Secret Manager**: https://cloud.google.com/secret-manager
- **Azure Key Vault**: https://azure.microsoft.com/en-us/services/key-vault/
- **HashiCorp Vault**: https://www.vaultproject.io/

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-05 | Security Team | Initial secrets management guide |

**Next Review Date:** 2026-02-05 (Quarterly)

---

**Classification:** Internal Use
**Distribution:** Engineering, DevOps, Security Teams

For questions or concerns, contact: security@univyn.com
