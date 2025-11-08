# Secure Deployment Checklist

**Version:** 1.0
**Last Updated:** 2025-11-05

This checklist ensures all security controls are properly configured before deploying to production.

---

## Pre-Deployment Checklist

### Environment Configuration

- [ ] **Environment Variable Set**
  - [ ] `ENVIRONMENT=production` configured
  - [ ] `SERVER_ADDRESS` set appropriately
  - [ ] `DATABASE_URL` configured with production credentials
  - [ ] `TOKEN_SYMMETRIC_KEY` set (≥32 random characters, unique to production)
  - [ ] All required environment variables present (check `.env.example`)

- [ ] **Secrets Management**
  - [ ] No secrets hardcoded in code
  - [ ] `.env` and `app.env` in `.gitignore`
  - [ ] Secrets stored in secure secret manager (AWS Secrets Manager, etc.)
  - [ ] Secrets rotation schedule documented
  - [ ] Emergency secret rotation procedure tested

- [ ] **Database Security**
  - [ ] `DATABASE_URL` uses SSL/TLS (`sslmode=require` or `sslmode=verify-full`)
  - [ ] Database credentials are strong (16+ characters)
  - [ ] Database credentials unique to production environment
  - [ ] Database user has least-privilege permissions
  - [ ] Database connection pool limits configured
  - [ ] Database backups enabled and encrypted
  - [ ] Backup restoration procedure tested

### Network Security

- [ ] **HTTPS/TLS**
  - [ ] Valid TLS certificate installed
  - [ ] TLS 1.2+ only (TLS 1.0/1.1 disabled)
  - [ ] Strong cipher suites configured
  - [ ] HTTP to HTTPS redirect enabled
  - [ ] HSTS header configured (will be set by application)

- [ ] **CORS Configuration**
  - [ ] `CORS_ALLOWED_ORIGINS` set to specific production domains
  - [ ] ❌ NOT using wildcard `*`
  - [ ] Domains verified and tested
  - [ ] Credentials allowed only for trusted origins

- [ ] **Firewall & Network**
  - [ ] Only necessary ports open (typically 443 for HTTPS)
  - [ ] Database not accessible from public internet
  - [ ] SSH access restricted to authorized IPs
  - [ ] Rate limiting configured at load balancer/WAF level

### Application Security

- [ ] **Rate Limiting**
  - [ ] `RATE_LIMIT_ENABLED=true`
  - [ ] `RATE_LIMIT_DEFAULT` configured appropriately (e.g., 100)
  - [ ] `RATE_LIMIT_AUTH` configured strictly (e.g., 5)
  - [ ] Rate limiting tested under load

- [ ] **Authentication & Authorization**
  - [ ] Token expiration configured appropriately (15 min access, 24h refresh)
  - [ ] RBAC middleware applied to all protected endpoints
  - [ ] No endpoints accessible without proper authorization
  - [ ] Account lockout protection enabled (migration applied)

- [ ] **Security Headers**
  - [ ] CSP header configured
  - [ ] HSTS header configured (auto-applied in production)
  - [ ] All security headers verified (see test output)

- [ ] **Request Protection**
  - [ ] Request size limits configured
  - [ ] Input validation enabled for all endpoints
  - [ ] File upload restrictions (if applicable)

### Code Quality & Security

- [ ] **Dependencies**
  - [ ] All dependencies updated to latest secure versions
  - [ ] No known high/critical vulnerabilities (`govulncheck ./...`)
  - [ ] Security scan passed (`gosec ./...`)
  - [ ] Dependency versions pinned in `go.mod`

- [ ] **Code Review**
  - [ ] Security-focused code review completed
  - [ ] No hardcoded credentials found
  - [ ] All database queries use SQLC (parameterized)
  - [ ] Input validation present on all endpoints
  - [ ] Proper error handling (no sensitive info in errors)
  - [ ] PII redacted in logs

- [ ] **Testing**
  - [ ] All unit tests passing
  - [ ] Integration tests passing
  - [ ] Security tests passing
  - [ ] Load testing completed
  - [ ] Penetration testing completed (if required)

### Database

- [ ] **Migrations**
  - [ ] All migrations applied successfully
  - [ ] Migration rollback tested
  - [ ] Database indexes created
  - [ ] Database constraints verified

- [ ] **Security Migrations**
  - [ ] Login attempts tracking tables created (migration 011)
  - [ ] Account lockout tables created
  - [ ] User security fields added (is_locked, failed_login_attempts, etc.)

### Logging & Monitoring

- [ ] **Logging**
  - [ ] Structured logging configured (zerolog)
  - [ ] Log level set appropriately (Info for production)
  - [ ] PII redaction verified
  - [ ] Sensitive data not logged (passwords, tokens, etc.)
  - [ ] Logs shipped to centralized logging system

- [ ] **Monitoring**
  - [ ] Application health check endpoint working (`/health`)
  - [ ] Monitoring dashboard created
  - [ ] Security metrics tracked:
    - [ ] Failed login attempts
    - [ ] Account lockouts
    - [ ] Authorization failures
    - [ ] Rate limit violations
    - [ ] Invalid token attempts
  - [ ] Alerting rules configured (see doc/SECURITY.md)

- [ ] **Audit Logging**
  - [ ] Security events logged:
    - [ ] Authentication events
    - [ ] Authorization failures
    - [ ] Privilege changes
    - [ ] Account lockouts
    - [ ] Administrative actions
  - [ ] Audit log retention configured (90+ days)

### Backup & Recovery

- [ ] **Backups**
  - [ ] Database backups enabled
  - [ ] Backup frequency configured (daily recommended)
  - [ ] Backup encryption enabled
  - [ ] Backup retention policy configured
  - [ ] Backup restoration tested successfully

- [ ] **Disaster Recovery**
  - [ ] Recovery Time Objective (RTO) documented
  - [ ] Recovery Point Objective (RPO) documented
  - [ ] Disaster recovery procedure documented
  - [ ] DR test completed

### Compliance & Documentation

- [ ] **Documentation**
  - [ ] Security documentation updated (doc/SECURITY.md)
  - [ ] Incident response procedures documented (doc/INCIDENT_RESPONSE.md)
  - [ ] Secrets management guide available (doc/SECRETS_MANAGEMENT.md)
  - [ ] Deployment runbook updated
  - [ ] API documentation updated

- [ ] **Compliance**
  - [ ] GDPR compliance verified (if applicable)
  - [ ] Data retention policies implemented
  - [ ] Privacy policy updated
  - [ ] Terms of service updated
  - [ ] Cookie policy updated (if web interface)

### Operational Readiness

- [ ] **Incident Response**
  - [ ] Incident response team identified
  - [ ] Emergency contacts documented
  - [ ] Incident response playbook accessible
  - [ ] Communication channels established
  - [ ] Incident response drill completed

- [ ] **On-Call**
  - [ ] On-call schedule established
  - [ ] On-call runbooks available
  - [ ] Escalation procedures documented
  - [ ] Alerting tested

---

## Deployment Steps

### Step 1: Pre-Deployment Verification

```bash
# 1. Verify all environment variables
env | grep -E "(DATABASE_URL|TOKEN_SYMMETRIC_KEY|CORS_ALLOWED_ORIGINS)"

# 2. Run security checks
gosec ./...
govulncheck ./...

# 3. Run all tests
make test

# 4. Build application
make build

# 5. Verify build artifacts
./bin/connect --version
```

### Step 2: Database Migration

```bash
# 1. Backup database
pg_dump $DATABASE_URL > backup-$(date +%Y%m%d-%H%M%S).sql

# 2. Run migrations (with transaction safety)
make migrate-up

# 3. Verify migration success
psql $DATABASE_URL -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"

# 4. Verify new tables exist
psql $DATABASE_URL -c "\dt login_attempts"
psql $DATABASE_URL -c "\dt account_lockouts"
```

### Step 3: Application Deployment

**Option A: Rolling Update (Zero Downtime)**
```bash
# Deploy to subset of servers first
kubectl set image deployment/connect-api connect-api=univyn/connect-api:v1.1.0

# Monitor health
watch kubectl get pods

# If successful, continue rollout
kubectl rollout status deployment/connect-api
```

**Option B: Blue-Green Deployment**
```bash
# Deploy to "green" environment
# Test thoroughly
# Switch traffic to green
# Keep blue as rollback option
```

### Step 4: Post-Deployment Verification

```bash
# 1. Health check
curl https://api.univyn.com/health

# 2. Test authentication
curl -X POST https://api.univyn.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'

# 3. Verify security headers
curl -I https://api.univyn.com/health | grep -E "(Strict-Transport-Security|Content-Security-Policy|X-Frame-Options)"

# 4. Test rate limiting
# Run 10 requests rapidly and verify rate limiting kicks in

# 5. Check logs
tail -f /var/log/connect/app.log

# 6. Monitor error rates
# Check monitoring dashboard for anomalies
```

### Step 5: Security Verification

```bash
# 1. Verify CORS configuration
curl -H "Origin: https://malicious.com" https://api.univyn.com/api/users
# Should be rejected

# 2. Verify HTTPS redirect
curl -I http://api.univyn.com/health
# Should redirect to HTTPS

# 3. Test account lockout
# Make 6 failed login attempts
# Verify account is locked

# 4. Verify authorization
# Try accessing admin endpoint without admin role
# Should return 403 Forbidden
```

---

## Rollback Procedure

If issues are discovered after deployment:

### Immediate Rollback

```bash
# 1. Rollback application
kubectl rollout undo deployment/connect-api

# 2. If database migrations were applied, rollback (if safe)
make migrate-down1

# 3. Verify rollback success
curl https://api.univyn.com/health

# 4. Investigate issue
# Review logs, monitor metrics

# 5. Fix issue in staging
# Re-test thoroughly

# 6. Plan next deployment
```

### Database Rollback (Use with Caution)

```bash
# Only if migrations caused issues and rollback is safe

# 1. Stop application
kubectl scale deployment/connect-api --replicas=0

# 2. Rollback migration
make migrate-down1

# 3. Verify database state
psql $DATABASE_URL -c "SELECT version FROM schema_migrations;"

# 4. Restart application with previous version
kubectl scale deployment/connect-api --replicas=3
```

---

## Post-Deployment Tasks

### Immediate (within 1 hour)

- [ ] Monitor error rates
- [ ] Monitor response times
- [ ] Monitor security events
- [ ] Verify no unexpected behavior
- [ ] Announce deployment in team channel

### Short-term (within 24 hours)

- [ ] Review logs for errors
- [ ] Check security alerts
- [ ] Verify backups running
- [ ] Monitor user feedback
- [ ] Document any issues encountered

### Medium-term (within 1 week)

- [ ] Review security metrics
- [ ] Analyze login attempt patterns
- [ ] Check for failed authorization attempts
- [ ] Review account lockout events
- [ ] Update documentation with lessons learned

---

## Security Verification Report

After deployment, generate and review security verification report:

### Required Checks

1. **Authentication**
   - [ ] Token verification working
   - [ ] Password validation enforced
   - [ ] Account lockout functioning

2. **Authorization**
   - [ ] RBAC checks enforced
   - [ ] Admin endpoints protected
   - [ ] Resource ownership respected

3. **Network Security**
   - [ ] HTTPS enforced
   - [ ] CORS configured correctly
   - [ ] Security headers present
   - [ ] Rate limiting active

4. **Data Protection**
   - [ ] Database encryption enabled
   - [ ] Passwords hashed with bcrypt
   - [ ] PII redacted in logs
   - [ ] SQL injection prevented

5. **Operational Security**
   - [ ] Monitoring active
   - [ ] Alerts configured
   - [ ] Logs shipping correctly
   - [ ] Backups running

### Verification Sign-off

| Area | Verified By | Date | Status |
|------|-------------|------|--------|
| Application Security | | | ⬜ Pass ⬜ Fail |
| Network Security | | | ⬜ Pass ⬜ Fail |
| Data Protection | | | ⬜ Pass ⬜ Fail |
| Operational Security | | | ⬜ Pass ⬜ Fail |

**Deployment approved for production:** ⬜ Yes ⬜ No

**Approver:** ___________________  **Date:** ___________

---

## Continuous Security

Post-deployment, maintain security posture:

### Daily

- [ ] Review security alerts
- [ ] Monitor failed login attempts
- [ ] Check for rate limit violations

### Weekly

- [ ] Review audit logs
- [ ] Check for locked accounts
- [ ] Review authorization failures
- [ ] Analyze suspicious patterns

### Monthly

- [ ] Update dependencies
- [ ] Run security scans
- [ ] Review access controls
- [ ] Test backup restoration
- [ ] Security metrics review

### Quarterly

- [ ] Rotate secrets
- [ ] Security team review
- [ ] Update threat model
- [ ] Conduct security drill
- [ ] Review and update documentation
- [ ] Penetration testing (if scheduled)

---

## Support & Resources

- **Security Documentation:** `doc/SECURITY.md`
- **Incident Response:** `doc/INCIDENT_RESPONSE.md`
- **Secrets Management:** `doc/SECRETS_MANAGEMENT.md`
- **Threat Model:** `doc/THREAT_MODEL.md`

- **Security Team:** security@univyn.com
- **On-Call:** [Phone Number]
- **Emergency Contacts:** See `doc/INCIDENT_RESPONSE.md`

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-05 | Security Team | Initial deployment checklist |

**Next Review Date:** 2026-02-05 (Quarterly)

---

**Print this checklist and check off each item during deployment.**
**Store completed checklists as deployment records.**
