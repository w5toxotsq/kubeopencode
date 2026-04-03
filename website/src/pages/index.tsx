import type {ReactNode} from 'react';
import {useEffect, useState} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import CodeBlock from '@theme/CodeBlock';

import styles from './index.module.css';

const GITHUB_REPO = 'kubeopencode/kubeopencode';

const agentYaml = `apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: dev-agent
spec:
  profile: "Interactive development agent"
  workspaceDir: /workspace
  port: 4096
  persistence:
    sessions:
      size: "2Gi"
  standby:
    idleTimeout: "30m"`;

const taskYaml = `apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: update-dependencies
spec:
  templateRef:
    name: ci-runner
  description: |
    Update all dependencies to latest versions.
    Run tests and create a pull request.`;

const helmInstallCmd = `helm repo add kubeopencode https://kubeopencode.github.io/kubeopencode
helm install kubeopencode kubeopencode/kubeopencode \\
  --namespace kubeopencode-system --create-namespace`;

// Hook to fetch GitHub star count
function useGitHubStars(): number | null {
  const [stars, setStars] = useState<number | null>(null);
  useEffect(() => {
    fetch(`https://api.github.com/repos/${GITHUB_REPO}`)
      .then(res => res.json())
      .then(data => {
        if (typeof data.stargazers_count === 'number') {
          setStars(data.stargazers_count);
        }
      })
      .catch(() => {});
  }, []);
  return stars;
}

function StarIcon({size = 16}: {size?: number}): ReactNode {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" fill="currentColor" style={{marginRight: '0.4rem', verticalAlign: 'text-bottom'}}>
      <path d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"/>
    </svg>
  );
}

// Section 1: Alpha Banner
function AlphaBanner(): ReactNode {
  return (
    <div className={styles.alphaBanner}>
      <div className="container">
        This project is in <strong>early alpha</strong> (v0.0.x). Not recommended for production use. API may change without backward compatibility.
      </div>
    </div>
  );
}

// Section 2: Hero
function HeroSection(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  const stars = useGitHubStars();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className={styles.heroTitle}>
          {siteConfig.title}
        </Heading>
        <p className={styles.heroTagline}>{siteConfig.tagline}</p>
        <p className={styles.heroDescription}>
          Deploy, manage, and scale AI coding agents on Kubernetes.
          Built on <a href="https://opencode.ai" className={styles.heroLink}>OpenCode</a>, designed for teams and enterprise.
        </p>
        <div className={styles.heroActions}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/getting-started">
            Get Started
          </Link>
          <Link
            className={styles.starBadge}
            href={`https://github.com/${GITHUB_REPO}`}>
            <span className={styles.starBadgeLeft}>
              <StarIcon size={16} />
              Star
            </span>
            {stars !== null && (
              <span className={styles.starBadgeCount}>{stars}</span>
            )}
          </Link>
        </div>
        <div className={styles.installSnippet}>
          <CodeBlock language="bash" title="Quick Install">
            {helmInstallCmd}
          </CodeBlock>
        </div>
      </div>
    </header>
  );
}

// Section 3: Demo placeholder
function DemoSection(): ReactNode {
  return (
    <section id="demo" className={styles.section}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">See KubeOpenCode in Action</Heading>
          <p className={styles.sectionSubtitle}>
            Watch how AI agents run natively on Kubernetes with zero configuration overhead.
          </p>
        </div>
        <div className={styles.demoPlaceholder}>
          <div className={styles.demoPlaceholderInner}>
            <div className={styles.demoPlayIcon}>
              <svg width="48" height="48" viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
            </div>
            <p className={styles.demoPlaceholderText}>Demo coming soon</p>
          </div>
        </div>
      </div>
    </section>
  );
}

// Section 4: Features
type FeatureItem = {
  title: string;
  icon: string;
  color: string;
  description: string;
};

const features: FeatureItem[] = [
  {
    title: 'Live Agents',
    icon: '\u26A1',
    color: '#f59e0b',
    description:
      'Every Agent runs as a persistent service on Kubernetes. Interactive terminal access, shared context across tasks, zero cold start \u2014 perfect for team-shared coding assistants.',
  },
  {
    title: 'Kubernetes-Native',
    icon: '\u2638\uFE0F',
    color: '#326ce5',
    description:
      'Declarative CRDs, GitOps-friendly, works with Helm, Kustomize, and ArgoCD. No new tools to learn \u2014 just kubectl apply.',
  },
  {
    title: 'Enterprise Ready',
    icon: '\uD83C\uDFE2',
    color: '#10b981',
    description:
      'RBAC, private registries, corporate proxies, custom CA certificates, pod security policies, and audit-ready infrastructure for compliance.',
  },
  {
    title: 'Built for Teams',
    icon: '\uD83D\uDC65',
    color: '#8b5cf6',
    description:
      'Shared agent configurations, batch operations, concurrency control, and centralized credential management for consistent standards.',
  },
];

function FeaturesSection(): ReactNode {
  return (
    <section id="features" className={clsx(styles.section, styles.sectionAlt)}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">Features</Heading>
          <p className={styles.sectionSubtitle}>
            Everything you need to run AI agents at scale on Kubernetes.
          </p>
        </div>
        <div className={styles.featuresGrid}>
          {features.map((feature, idx) => (
            <div key={idx} className={styles.featureCard}>
              <div className={styles.featureIconArea} style={{backgroundColor: feature.color}}>
                <span className={styles.featureIcon}>{feature.icon}</span>
              </div>
              <Heading as="h3" className={styles.featureTitle}>{feature.title}</Heading>
              <p className={styles.featureDescription}>{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// Section 5: How It Works (refactored QuickExample)
function HowItWorksSection(): ReactNode {
  return (
    <section id="how-it-works" className={styles.section}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">How It Works</Heading>
          <p className={styles.sectionSubtitle}>
            Two simple resources. One powerful workflow.
          </p>
        </div>
        <div className={styles.stepsContainer}>
          <div className={styles.step}>
            <div className={styles.stepNumber}>1</div>
            <div className={styles.stepContent}>
              <Heading as="h3">Define an Agent</Heading>
              <p>
                Deploy persistent AI agents your team can interact with in real time
                &mdash; through the web terminal, CLI, or by submitting Tasks.
              </p>
              <ul className={styles.stepList}>
                <li>Zero cold start &mdash; agent is always running</li>
                <li>Interactive terminal access via CLI or web</li>
                <li>Auto-suspend when idle, resume on demand</li>
                <li>Session history persists across restarts</li>
              </ul>
            </div>
            <div className={styles.stepCode}>
              <CodeBlock language="yaml" title="agent.yaml">
                {agentYaml}
              </CodeBlock>
            </div>
          </div>
          <div className={styles.step}>
            <div className={styles.stepNumber}>2</div>
            <div className={styles.stepContent}>
              <Heading as="h3">Submit a Task</Heading>
              <p>
                Run stable, repeatable AI tasks in ephemeral Pods. Perfect for
                CI/CD pipelines, batch operations, and automated workflows.
              </p>
              <ul className={styles.stepList}>
                <li>No new tools to learn &mdash; just <code>kubectl apply</code></li>
                <li>Works with any CI/CD pipeline</li>
                <li>Scale with Helm templates for batch operations</li>
                <li>Rate limiting and quota controls</li>
              </ul>
            </div>
            <div className={styles.stepCode}>
              <CodeBlock language="yaml" title="task.yaml">
                {taskYaml}
              </CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// Section 6: Architecture
function ArchitectureSection(): ReactNode {
  return (
    <section id="architecture" className={clsx(styles.section, styles.sectionAlt)}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">Architecture</Heading>
          <p className={styles.sectionSubtitle}>
            A simple, Kubernetes-native design with no external dependencies.
          </p>
        </div>
        <div className={styles.archFlow}>
          <div className={styles.archNode}>
            <div className={styles.archNodeIcon}>
              <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm-1 2l5 5h-5V4zM6 20V4h5v7h7v9H6z"/>
              </svg>
            </div>
            <Heading as="h4">Task</Heading>
            <p>WHAT to do</p>
          </div>
          <div className={styles.archArrow}>
            <svg width="48" height="24" viewBox="0 0 48 24" fill="currentColor">
              <path d="M0 11h40l-6-6 1.4-1.4L44 12l-8.6 8.4L34 19l6-6H0z"/>
            </svg>
          </div>
          <div className={styles.archNode}>
            <div className={styles.archNodeIcon}>
              <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
            <Heading as="h4">Agent</Heading>
            <p>HOW to execute</p>
          </div>
          <div className={styles.archArrow}>
            <svg width="48" height="24" viewBox="0 0 48 24" fill="currentColor">
              <path d="M0 11h40l-6-6 1.4-1.4L44 12l-8.6 8.4L34 19l6-6H0z"/>
            </svg>
          </div>
          <div className={styles.archNode}>
            <div className={styles.archNodeIcon}>
              <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20 3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H4V5h16v14zM6 15h2v2H6zm0-4h2v2H6zm0-4h2v2H6zm4 8h8v2h-8zm0-4h8v2h-8zm0-4h8v2h-8z"/>
              </svg>
            </div>
            <Heading as="h4">Pod (OpenCode)</Heading>
            <p>WHERE it runs</p>
          </div>
        </div>
        <div className={styles.archDetails}>
          <div className={styles.archDetail}>
            <strong>No external dependencies</strong>
            <span>No PostgreSQL, no Redis &mdash; just Kubernetes (etcd for state, Pods for execution)</span>
          </div>
          <div className={styles.archDetail}>
            <strong>Single container image</strong>
            <span>Controller, init containers, and utilities all in one image</span>
          </div>
          <div className={styles.archDetail}>
            <strong>Two-container pattern</strong>
            <span>Init container copies OpenCode binary, worker container runs the agent</span>
          </div>
        </div>
      </div>
    </section>
  );
}

// Section 7: FAQ
type FaqItem = {
  question: string;
  answer: string;
};

const faqItems: FaqItem[] = [
  {
    question: 'What is KubeOpenCode?',
    answer:
      'KubeOpenCode is an open-source platform that brings AI coding agents into the Kubernetes ecosystem. It lets you deploy, manage, and scale AI agents as Kubernetes-native services with declarative CRDs, RBAC, and enterprise-grade governance.',
  },
  {
    question: 'How is KubeOpenCode different from OpenCode?',
    answer:
      'OpenCode is a terminal-based AI coding tool for individual developers. KubeOpenCode builds on top of OpenCode to provide a Kubernetes-native platform for teams and enterprise \u2014 adding multi-tenancy, persistent agents, task orchestration, and infrastructure-level controls.',
  },
  {
    question: 'Do I need to run my own Kubernetes cluster?',
    answer:
      'Yes, KubeOpenCode runs on any standard Kubernetes cluster (v1.26+). You can use managed services like EKS, GKE, AKS, or a local cluster with Kind or minikube for development.',
  },
  {
    question: 'What AI models are supported?',
    answer:
      'KubeOpenCode supports any model that OpenCode supports, including Claude, GPT-4, Gemini, and local models via Ollama. You configure model access through Agent credentials and the OpenCode configuration.',
  },
  {
    question: 'Is KubeOpenCode free?',
    answer:
      'Yes, KubeOpenCode is free and open-source under the Apache License 2.0. You only pay for your own infrastructure (Kubernetes cluster) and AI model API usage.',
  },
  {
    question: 'Can I use it in production?',
    answer:
      'KubeOpenCode is currently in early alpha (v0.0.x). The API may change without backward compatibility. We recommend using it for development, testing, and evaluation while we stabilize the API toward a v1.0 release.',
  },
  {
    question: 'How do I contribute?',
    answer:
      'Contributions are welcome! Check out our GitHub repository for open issues, join the Slack community for discussions, or submit a pull request. See the CONTRIBUTING guide in the repository for details.',
  },
];

function FaqSection(): ReactNode {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleFaq = (index: number): void => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <section id="faq" className={styles.section}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">Frequently Asked Questions</Heading>
        </div>
        <div className={styles.faqList}>
          {faqItems.map((item, idx) => (
            <div
              key={idx}
              className={clsx(styles.faqItem, openIndex === idx && styles.faqItemOpen)}
            >
              <button
                className={styles.faqQuestion}
                onClick={() => toggleFaq(idx)}
                type="button"
                aria-expanded={openIndex === idx}
              >
                <span>{item.question}</span>
                <span className={styles.faqChevron}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M7.41 8.59L12 13.17l4.59-4.58L18 10l-6 6-6-6z"/>
                  </svg>
                </span>
              </button>
              <div className={styles.faqAnswer}>
                <p>{item.answer}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// Section 8: CTA
function CtaSection(): ReactNode {
  return (
    <section className={clsx(styles.section, styles.ctaSection)}>
      <div className="container">
        <div className={styles.ctaContent}>
          <Heading as="h2" className={styles.ctaTitle}>
            Ready to run AI agents on Kubernetes?
          </Heading>
          <p className={styles.ctaDescription}>
            Get started in minutes with Helm. Deploy your first agent today.
          </p>
          <div className={styles.ctaActions}>
            <Link
              className="button button--primary button--lg"
              to="/docs/getting-started">
              Get Started
            </Link>
            <Link
              className="button button--outline button--lg"
              href={`https://github.com/${GITHUB_REPO}`}>
              View on GitHub
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="Kubernetes-native Agent Platform for Teams and Enterprise"
      description="Run AI agents on Kubernetes. Built on OpenCode, designed for teams and enterprise.">
      <AlphaBanner />
      <HeroSection />
      <main>
        <DemoSection />
        <FeaturesSection />
        <HowItWorksSection />
        <ArchitectureSection />
        <FaqSection />
        <CtaSection />
      </main>
    </Layout>
  );
}
