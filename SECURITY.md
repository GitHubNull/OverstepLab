# Legal Notice and Security Disclaimer

## 法律声明与安全免责声明

### English

**OverstepLab is a deliberately vulnerable web application designed exclusively for authorized security education, training, and research.**

#### ⚠️ CRITICAL WARNINGS

1. **Intentional Vulnerabilities**: This software contains intentionally implanted security vulnerabilities including broken access control, insecure direct object references, and privilege escalation flaws.

2. **NOT FOR PRODUCTION USE**: It is NOT SECURE and must NEVER be deployed to production environments or exposed to the public internet.

3. **Isolated Environment Required**: It must ONLY be run in isolated, controlled environments (local machine, private VM, or closed container network) with no access to sensitive data.

4. **Legal Compliance**: DO NOT use the techniques demonstrated by this software against any system, network, or application without explicit written authorization from the owner.

#### LIABILITY

The authors, contributors, and maintainers of OverstepLab assume **NO LIABILITY** for any misuse, damage, data loss, or legal consequences arising from the use of this software. Users bear sole responsibility for ensuring their use complies with all applicable local, national, and international laws.

By downloading, building, or running OverstepLab, you acknowledge that you understand these risks and agree to use this software solely for lawful educational purposes.

---

### 中文

**OverstepLab 是一个故意包含安全漏洞的 Web 应用程序，仅供授权的安全教育、培训和研究使用。**

#### ⚠️ 重要警告

1. **故意漏洞**：本软件包含故意植入的安全漏洞，包括但不限于越权访问、不安全的直接对象引用和权限提升缺陷。

2. **禁止生产使用**：本软件不安全，绝对禁止部署到生产环境或暴露在公网。

3. **需要隔离环境**：本软件只能在隔离的、可控的环境中运行（本地机器、私有虚拟机或封闭容器），且不得接触敏感数据。

4. **合法合规**：未经所有者明确书面授权，禁止使用本软件演示的技术攻击任何系统、网络或应用。

#### 责任免除

OverstepLab 的作者、贡献者和维护者对本软件的使用产生的任何滥用、损害、数据丢失或法律后果承担**零责任**。用户须自行确保其使用符合所有适用的当地、国家和国际法律。

下载、构建或运行 OverstepLab 即表示您已知晓这些风险，并同意仅将本软件用于合法的教育目的。

---

## Vulnerability List

The following vulnerabilities are intentionally implanted:

| ID | Category | Description |
|----|----------|-------------|
| H-01~H-05 | Horizontal IDOR | Insecure direct object reference allowing access to other users' resources |
| V-01~V-05 | Vertical Escalation | Missing role/permission checks allowing low-privilege users to perform admin actions |
| C-01~C-03 | Context Escalation | Missing business context validation allowing cross-boundary access |
