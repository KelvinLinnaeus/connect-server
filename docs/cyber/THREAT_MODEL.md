# Threat Model - Univyn Connect Server

**Version:** 1.0
**Date:** 2025-11-05
**Classification:** Internal Use

## Executive Summary

This document provides a comprehensive threat model for the Univyn Connect Server, a multi-tenant social networking backend built with Go and PostgreSQL. It identifies threat agents, attack surfaces, sensitive assets, high-risk scenarios, and prioritized security controls.

---

## 1. System Overview

### Architecture
- **Application:** Go 1.24.0 with Gin web framework
- **Database:** PostgreSQL 15+
- **Authentication:** PASETO v2 symmetric tokens
- **Deployment:** Cloud-hosted (AWS/Neon)
- **API:** RESTful JSON over HTTPS
- **Multi-tenancy:** Space-based isolation

### Key Components
1. **API Layer** - 80+ HTTP endpoints
2. **Authentication Service** - Login, registration, token management
3. **Authorization Service** - Role-based access control (RBAC)
4. **Business Logic** - 14 service packages
5. **Data Layer** - PostgreSQL with SQLC
6. **Session Management** - Database-backed refresh tokens

---

## 2. Threat Agents

### External Threat Agents

| Agent | Motivation | Capability | Likelihood |
|-------|------------|------------|------------|
| **Opportunistic Attackers** | Financial gain, defacement | Low-Medium | High |
| **Credential Stuffers** | Account takeover, data theft | Medium | High |
| **API Abusers** | Resource exhaustion, scraping | Low-Medium | Medium |
| **Social Engineers** | Account compromise | Low | Medium |
| **Competitors** | Intelligence gathering | Medium-High | Low |
| **Advanced Persistent Threats (APT)** | Long-term access, espionage | High | Very Low |

### Internal Threat Agents

| Agent | Motivation | Capability | Likelihood |
|-------|------------|------------|------------|
| **Malicious Insiders** | Data theft, sabotage | High | Low |
| **Compromised Accounts** | Lateral movement | Medium | Medium |
| **Negligent Employees** | Accidental exposure | Low | Medium |

### Automated Threat Agents

| Agent | Motivation | Capability | Likelihood |
|-------|------------|------------|------------|
| **Bots & Scrapers** | Data harvesting | Low-Medium | High |
| **Vulnerability Scanners** | Reconnaissance | Medium | High |
| **DDoS Networks** | Service disruption | High | Medium |

---

## 3. Assets & Data Classification

### Critical Assets

| Asset | Confidentiality | Integrity | Availability | Impact if Compromised |
|-------|----------------|-----------|--------------|----------------------|
| **User Credentials** | CRITICAL | CRITICAL | HIGH | Account takeover, identity theft |
| **Authentication Tokens** | CRITICAL | CRITICAL | HIGH | Unauthorized access |
| **Personal Identifiable Information (PII)** | CRITICAL | HIGH | MEDIUM | Privacy violations, GDPR fines |
| **Private Messages** | CRITICAL | HIGH | MEDIUM | Privacy breach, reputational damage |
| **Database Credentials** | CRITICAL | CRITICAL | CRITICAL | Complete system compromise |
| **API Keys** | HIGH | HIGH | HIGH | Service abuse |
| **User Content** | MEDIUM | HIGH | MEDIUM | Data loss, defacement |
| **Session Data** | HIGH | HIGH | MEDIUM | Session hijacking |

### Sensitive Operations

1. **Authentication & Authorization**
   - Login/logout
   - Password reset
   - Token issuance
   - Role assignment

2. **Data Modification**
   - User profile updates
   - Content creation/deletion
   - Relationship changes (follow, block)

3. **Administrative Functions**
   - User management
   - Space/community management
   - Moderation actions

4. **Financial Transactions** (if applicable)
   - Payments
   - Subscriptions

---

## 4. Attack Surface Inventory

### Network Attack Surfaces

#### 4.1 HTTP API Endpoints

**Public Endpoints (Unauthenticated)**
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - Authentication
- `POST /api/auth/forgot-password` - Password reset initiation
- `POST /api/auth/reset-password` - Password reset completion
- `GET /api/health` - Health check
- `GET /api/spaces` - Public space listing

**Risk:** Brute force, credential stuffing, enumeration, abuse

**Protected Endpoints (Authenticated - 70+ endpoints)**
- User profile operations
- Social interactions (posts, comments, likes)
- Communities & groups
- Direct messaging
- Events & registrations
- Notifications

**Risk:** Unauthorized access, privilege escalation, data exfiltration

**Administrative Endpoints**
- Space management
- User administration
- Content moderation

**Risk:** Privilege escalation, administrative abuse

#### 4.2 Authentication Mechanisms

| Mechanism | Attack Vector | Mitigation Status |
|-----------|---------------|-------------------|
| Password Authentication | Brute force, credential stuffing | ⚠️ Partial (rate limiting exists) |
| Token-based Auth (PASETO) | Token theft, replay | ✅ Strong encryption |
| Refresh Tokens | Token theft | ⚠️ Limited revocation |
| Session Management | Session hijacking | ⚠️ Needs enhancement |

#### 4.3 Database Interface

- **PostgreSQL Connection**: Direct access from application
- **Connection String**: Contains credentials
- **Query Interface**: SQLC-generated prepared statements

**Risk:** SQL injection (mitigated), credential exposure (CRITICAL), connection pool exhaustion

### Application Attack Surfaces

#### 4.4 Input Validation

**User-Controlled Inputs:**
- JSON request bodies
- Query parameters
- Path parameters
- HTTP headers
- File uploads (if present)

**Risk:** Injection attacks, buffer overflow, type confusion, malicious file upload

#### 4.5 Business Logic

**Complex Operations:**
- Multi-step workflows (registration, password reset)
- State transitions (event registration, follow/unfollow)
- Cascading deletions
- Notification triggers

**Risk:** Race conditions, state manipulation, authorization bypass

### Infrastructure Attack Surfaces

#### 4.6 Configuration

**Sensitive Configuration:**
- Database credentials
- Token signing keys
- API keys for third-party services
- CORS origins
- TLS certificates

**Risk:** Credential exposure, misconfiguration

#### 4.7 Dependencies

**External Dependencies:**
- Go modules (50+ packages)
- Database drivers
- HTTP framework (Gin)
- Logging libraries

**Risk:** Supply chain attacks, known vulnerabilities

#### 4.8 Deployment & Runtime

**Infrastructure:**
- Cloud hosting platform
- Database service (Neon)
- CI/CD pipeline
- Logging & monitoring

**Risk:** Cloud misconfigurations, compromised CI/CD, log injection

---

## 5. Threat Scenarios (Prioritized)

### 5.1 CRITICAL Threats

#### T1: Database Credential Exposure
**Description:** Hardcoded credentials in Makefile and test files
**Attack Vector:** Source code access (git clone, leaked repo)
**Impact:** Complete database compromise, data breach
**Likelihood:** HIGH (credentials already in repo)
**CVSS Score:** 10.0 (Critical)
**Mitigation Priority:** IMMEDIATE

#### T2: Authorization Bypass
**Description:** Incomplete RBAC implementation allows unauthorized access to privileged operations
**Attack Vector:** Direct API calls to admin endpoints with valid but non-privileged token
**Impact:** Privilege escalation, unauthorized data access/modification
**Likelihood:** HIGH
**CVSS Score:** 9.1 (Critical)
**Mitigation Priority:** IMMEDIATE

#### T3: Token Theft & Replay
**Description:** Stolen access tokens can be used without detection
**Attack Vector:** XSS in client, man-in-the-middle, malware
**Impact:** Account takeover, unauthorized actions
**Likelihood:** MEDIUM
**CVSS Score:** 8.1 (High)
**Mitigation Priority:** HIGH

### 5.2 HIGH Threats

#### T4: Credential Stuffing Attack
**Description:** Automated attempts using leaked credentials from other breaches
**Attack Vector:** Bulk login requests to `/api/auth/login`
**Impact:** Multiple account takeovers
**Likelihood:** HIGH
**CVSS Score:** 7.5 (High)
**Mitigation Priority:** HIGH

#### T5: API Abuse & Resource Exhaustion
**Description:** Excessive requests overwhelm server resources
**Attack Vector:** Automated scripts, botnets
**Impact:** Service degradation, downtime
**Likelihood:** MEDIUM
**CVSS Score:** 7.1 (High)
**Mitigation Priority:** HIGH

#### T6: Cross-Site Request Forgery (CSRF)
**Description:** Malicious site tricks authenticated user into making unwanted requests
**Attack Vector:** Social engineering with embedded forms/scripts
**Impact:** Unauthorized actions on behalf of victim
**Likelihood:** MEDIUM
**CVSS Score:** 6.5 (Medium)
**Mitigation Priority:** MEDIUM

#### T7: CORS Misconfiguration
**Description:** Wildcard CORS allows any origin to make authenticated requests
**Attack Vector:** Malicious JavaScript from attacker-controlled domain
**Impact:** Cross-origin data exfiltration
**Likelihood:** MEDIUM (if default config used)
**CVSS Score:** 6.8 (Medium)
**Mitigation Priority:** MEDIUM

### 5.3 MEDIUM Threats

#### T8: Information Disclosure via Error Messages
**Description:** Detailed error messages leak system information
**Attack Vector:** Invalid requests trigger verbose errors
**Impact:** Reconnaissance for further attacks
**Likelihood:** LOW
**CVSS Score:** 4.3 (Medium)
**Mitigation Priority:** MEDIUM

#### T9: Injection Attacks (XSS, SQL)
**Description:** Malicious input executed as code
**Attack Vector:** Unsanitized user input
**Impact:** XSS in clients, SQL injection (unlikely due to SQLC)
**Likelihood:** LOW (SQLC protects SQL, but XSS possible in clients)
**CVSS Score:** 6.1 (Medium)
**Mitigation Priority:** MEDIUM

#### T10: Dependency Vulnerabilities
**Description:** Outdated dependencies contain known CVEs
**Attack Vector:** Exploiting public vulnerabilities
**Impact:** Varies (RCE, DoS, information disclosure)
**Likelihood:** MEDIUM
**CVSS Score:** Varies
**Mitigation Priority:** MEDIUM

### 5.4 LOW Threats

#### T11: Session Fixation
**Description:** Attacker forces victim to use known session ID
**Attack Vector:** Social engineering, URL manipulation
**Impact:** Session hijacking
**Likelihood:** LOW (modern token implementation)
**CVSS Score:** 4.6 (Medium)
**Mitigation Priority:** LOW

#### T12: Directory Traversal in File Operations
**Description:** Malicious file paths access unauthorized files
**Attack Vector:** Crafted file upload/download requests
**Impact:** Unauthorized file access
**Likelihood:** LOW (file operations not evident in code review)
**CVSS Score:** 5.3 (Medium)
**Mitigation Priority:** LOW (if file operations exist)

---

## 6. Data Flow Analysis

### Authentication Flow

```
Client → POST /api/auth/login → Auth Handler → Verify Password (bcrypt)
  ↓
  → Generate PASETO Token (15 min) + Refresh Token (24h) → Store Session in DB
  ↓
  → Return Tokens to Client
```

**Threats:**
- T4: Credential stuffing at login endpoint
- T1: Database credential exposure
- T3: Token theft in transit (if HTTPS not enforced)

### Authorization Flow

```
Client → Request with Token → Auth Middleware → Verify PASETO → Extract User
  ↓
  → ⚠️ Authorization Middleware (INCOMPLETE) → Handler → Business Logic
  ↓
  → Database Query (SQLC) → Response
```

**Threats:**
- T2: Authorization bypass due to incomplete RBAC
- T3: Token replay
- T9: SQL injection (mitigated by SQLC)

### Data Storage Flow

```
Request → Validation (validator) → Service Layer → Database (PostgreSQL)
  ↓
  → Encrypted at Rest (if configured) → Backups
```

**Threats:**
- T1: Database credential exposure
- T9: Injection if validation bypassed
- T8: Information disclosure in error responses

---

## 7. Trust Boundaries

1. **Internet ↔ API Server**
   - Untrusted external requests
   - Must validate all inputs
   - Enforce authentication & rate limiting

2. **API Server ↔ Database**
   - Trusted but over network
   - Requires encrypted connections
   - Least-privilege credentials

3. **API Server ↔ External Services**
   - Third-party APIs (if any)
   - Mutual TLS preferred
   - API key rotation

4. **User ↔ User (Social Interactions)**
   - Trust but verify
   - Content moderation
   - Privacy controls

---

## 8. Security Controls (Current State)

### Implemented Controls ✅

| Control | Status | Coverage |
|---------|--------|----------|
| Password Hashing (bcrypt) | ✅ Strong | All passwords |
| PASETO Tokens | ✅ Strong | Authentication |
| SQL Injection Protection (SQLC) | ✅ Strong | All queries |
| Rate Limiting (IP-based) | ✅ Partial | Global (100 req/min) |
| Panic Recovery | ✅ Complete | All endpoints |
| Structured Error Responses | ✅ Good | All handlers |
| Input Validation Framework | ✅ Available | Inconsistent usage |

### Missing/Incomplete Controls ⚠️

| Control | Status | Risk |
|---------|--------|------|
| RBAC Authorization | ⚠️ Incomplete | CRITICAL |
| Token Revocation | ⚠️ Limited | HIGH |
| Account Lockout | ❌ Missing | HIGH |
| CSRF Protection | ❌ Missing | MEDIUM |
| Secure Headers | ❌ Missing | MEDIUM |
| Endpoint Rate Limiting | ⚠️ Partial | MEDIUM |
| Request Size Limits | ❌ Missing | MEDIUM |
| CORS Restrictions | ⚠️ Default Wildcard | MEDIUM |
| HTTPS Enforcement | ❌ Missing | HIGH |
| Audit Logging | ⚠️ Basic | MEDIUM |
| Secret Management | ❌ Hardcoded | CRITICAL |

---

## 9. Attack Trees

### Attack Tree 1: Account Takeover

```
[GOAL: Take over user account]
├── [A1] Credential Compromise
│   ├── [A1.1] Brute Force Login ⚠️ (Rate limiting exists but weak)
│   ├── [A1.2] Credential Stuffing ⚠️ (No account lockout)
│   ├── [A1.3] Phishing (Out of scope)
│   └── [A1.4] Password Reset Abuse ⚠️ (Needs review)
├── [A2] Token Theft
│   ├── [A2.1] Intercept in Transit ⚠️ (If HTTPS not enforced)
│   ├── [A2.2] XSS in Client (Client-side security)
│   └── [A2.3] Malware (Out of scope)
└── [A3] Session Hijacking
    ├── [A3.1] Session Fixation ✅ (PASETO resistant)
    └── [A3.2] Token Replay ⚠️ (No revocation)
```

### Attack Tree 2: Privilege Escalation

```
[GOAL: Gain administrative access]
├── [B1] Authorization Bypass
│   ├── [B1.1] Missing RBAC Checks ❌ CRITICAL
│   ├── [B1.2] Role Manipulation ⚠️ (Needs testing)
│   └── [B1.3] Parameter Tampering ⚠️ (Needs validation)
├── [B2] Compromise Admin Account
│   └── (See Attack Tree 1)
└── [B3] Exploit Business Logic
    ├── [B3.1] Race Conditions (Unknown)
    └── [B3.2] State Manipulation (Unknown)
```

### Attack Tree 3: Data Exfiltration

```
[GOAL: Steal user data]
├── [C1] Database Compromise
│   ├── [C1.1] Credential Exposure ❌ CRITICAL (Hardcoded)
│   └── [C1.2] SQL Injection ✅ (SQLC protected)
├── [C2] API Abuse
│   ├── [C2.1] Unauthorized Access ⚠️ (Authorization incomplete)
│   ├── [C2.2] Scraping ⚠️ (Rate limiting exists)
│   └── [C2.3] IDOR Attacks (Needs testing)
└── [C3] CORS Bypass
    └── [C3.1] Cross-Origin Requests ⚠️ (Wildcard default)
```

---

## 10. Risk Matrix

| Threat ID | Threat | Likelihood | Impact | Risk Level | Priority |
|-----------|--------|------------|--------|------------|----------|
| T1 | Database Credential Exposure | HIGH | CRITICAL | CRITICAL | IMMEDIATE |
| T2 | Authorization Bypass | HIGH | CRITICAL | CRITICAL | IMMEDIATE |
| T3 | Token Theft & Replay | MEDIUM | HIGH | HIGH | HIGH |
| T4 | Credential Stuffing | HIGH | HIGH | HIGH | HIGH |
| T5 | API Abuse | MEDIUM | HIGH | HIGH | HIGH |
| T6 | CSRF | MEDIUM | MEDIUM | MEDIUM | MEDIUM |
| T7 | CORS Misconfiguration | MEDIUM | MEDIUM | MEDIUM | MEDIUM |
| T8 | Information Disclosure | LOW | MEDIUM | MEDIUM | MEDIUM |
| T9 | Injection Attacks | LOW | MEDIUM | MEDIUM | MEDIUM |
| T10 | Dependency Vulnerabilities | MEDIUM | VARIES | MEDIUM | MEDIUM |
| T11 | Session Fixation | LOW | MEDIUM | LOW | LOW |
| T12 | Directory Traversal | LOW | MEDIUM | LOW | LOW |

---

## 11. Recommended Mitigations (Prioritized)

### Phase 1: IMMEDIATE (Critical Risks)

1. **Remove Hardcoded Credentials** (T1)
   - Move all credentials to environment variables
   - Remove from Makefile and test files
   - Update .gitignore to prevent future exposure
   - Rotate exposed credentials immediately

2. **Implement Complete RBAC** (T2)
   - Finish authorization middleware implementation
   - Add role checking for admin/moderator endpoints
   - Add authorization tests for all protected routes

### Phase 2: HIGH Priority (Within 1 Week)

3. **Account Protection** (T4)
   - Implement account lockout after N failed attempts
   - Add stricter rate limiting on auth endpoints (5 req/min)
   - Log suspicious authentication patterns

4. **Token Security** (T3)
   - Add token revocation capability
   - Implement token blacklist or versioning
   - Add security logging for token usage

5. **HTTPS Enforcement** (T3, T6)
   - Add middleware to redirect HTTP → HTTPS
   - Set Secure flag on all cookies
   - Enforce HSTS header

### Phase 3: MEDIUM Priority (Within 2 Weeks)

6. **Secure Headers** (T6, T7, T8)
   - Add CSP, HSTS, X-Frame-Options, X-Content-Type-Options
   - Fix CORS to whitelist specific origins
   - Add CSRF protection for state-changing operations

7. **Request Protection** (T5)
   - Add request body size limits
   - Implement endpoint-specific rate limits
   - Add connection limits

8. **Input Validation** (T9)
   - Audit all endpoints for input validation
   - Ensure consistent use of validator framework
   - Add output encoding where needed

9. **Audit Logging** (All Threats)
   - Log all authentication events
   - Log authorization failures
   - Log privilege changes
   - Ensure PII redaction

10. **Error Handling** (T8)
    - Review all error responses
    - Remove stack traces in production
    - Standardize error format

### Phase 4: ONGOING

11. **Dependency Management** (T10)
    - Update all dependencies to latest secure versions
    - Add automated vulnerability scanning (gosec, Snyk)
    - Pin dependency versions

12. **Security Testing** (All Threats)
    - Create comprehensive security test suite
    - Add CI security checks
    - Conduct penetration testing

13. **Monitoring & Response**
    - Implement security monitoring
    - Create incident response playbook
    - Set up alerting for suspicious activity

---

## 12. Assumptions & Constraints

### Assumptions
1. Application runs in trusted cloud environment
2. TLS termination handled by load balancer or reverse proxy
3. Database backups are encrypted and access-controlled
4. No direct database access from internet
5. CI/CD pipeline has secret management

### Constraints
1. Must maintain API backward compatibility
2. Cannot break existing client integrations
3. Zero-downtime deployment required
4. GDPR/privacy compliance required
5. Performance impact must be minimal

### Out of Scope
1. Client-side security (XSS, CSRF in web app)
2. Mobile app security
3. Physical security
4. Social engineering
5. DNS and network infrastructure
6. DDoS mitigation (assumed handled by CDN/WAF)

---

## 13. Compliance & Regulatory

### GDPR Requirements
- **Data Protection by Design**: Implement encryption, access controls
- **Right to Erasure**: Ensure data deletion is complete
- **Data Breach Notification**: 72-hour notification process
- **Privacy by Default**: Minimal data collection

### Industry Standards
- **OWASP Top 10**: Address all applicable risks
- **OWASP ASVS**: Follow Application Security Verification Standard
- **CWE Top 25**: Mitigate most dangerous software weaknesses
- **PCI DSS**: If handling payments (scope TBD)

---

## 14. Metrics & Monitoring

### Security Metrics to Track

1. **Authentication Metrics**
   - Failed login attempts per hour/day
   - Account lockouts triggered
   - Password reset requests
   - Suspicious login patterns (geolocation, time)

2. **Authorization Metrics**
   - Authorization failures per endpoint
   - Privilege escalation attempts
   - Admin action audit trail

3. **API Security Metrics**
   - Rate limit violations per endpoint
   - Request size violations
   - Invalid token attempts
   - CORS violations

4. **Application Security Metrics**
   - Input validation failures
   - Error rates per endpoint
   - Dependency vulnerabilities detected
   - Security test pass rates

### Alerting Thresholds

| Event | Threshold | Severity |
|-------|-----------|----------|
| Failed logins from single IP | >10 in 5 min | HIGH |
| Account lockouts | >5 in 1 hour | MEDIUM |
| Authorization failures | >50 in 5 min | HIGH |
| Rate limit violations | >100 in 1 min | MEDIUM |
| Invalid tokens | >20 in 5 min | MEDIUM |
| Database connection failures | >5 in 1 min | CRITICAL |

---

## 15. Review & Update Schedule

- **Quarterly Review**: Update threat model with new features/threats
- **Post-Incident Review**: Update after any security incident
- **Annual Penetration Test**: External security assessment
- **Continuous Monitoring**: Automated vulnerability scanning

---

## 16. References

1. **OWASP Resources**
   - OWASP Top 10: https://owasp.org/Top10/
   - OWASP ASVS: https://owasp.org/www-project-application-security-verification-standard/
   - OWASP Testing Guide: https://owasp.org/www-project-web-security-testing-guide/

2. **Standards**
   - NIST Cybersecurity Framework
   - CWE/SANS Top 25: https://cwe.mitre.org/top25/
   - CVSS Specification: https://www.first.org/cvss/

3. **Go Security**
   - Go Security Best Practices: https://go.dev/doc/security/
   - gosec - Go Security Checker: https://github.com/securego/gosec

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-05 | Security Team | Initial threat model |

**Next Review Date:** 2026-02-05 (Quarterly)

---

**Classification:** Internal Use
**Distribution:** Engineering, Security, Operations Teams

