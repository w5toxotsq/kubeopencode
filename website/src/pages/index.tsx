import type {ReactNode} from 'react';
import {useEffect, useState} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
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

function AlphaBanner() {
  return (
    <div className={styles.alphaBanner}>
      <div className="container">
        This project is in <strong>early alpha</strong> (v0.0.x). Not recommended for production use. API may change without backward compatibility.
      </div>
    </div>
  );
}

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  const stars = useGitHubStars();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <p className={styles.heroDescription}>
          Run AI coding agents as live services on Kubernetes. Deploy persistent agents
          your team can interact with anytime &mdash; or run batch tasks at scale.
          Built on OpenCode, designed for teams and enterprise.
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/getting-started">
            Get Started
          </Link>
          <Link
            className="button button--outline button--lg"
            style={{color: 'white', borderColor: 'white', marginLeft: '1rem'}}
            href={`https://github.com/${GITHUB_REPO}`}>
            GitHub{stars !== null ? ` (${stars})` : ''}
          </Link>
        </div>
      </div>
    </header>
  );
}

function QuickExample() {
  return (
    <section className={styles.quickExample}>
      <div className="container">
        <div className="row">
          <div className="col col--6">
            <Heading as="h2">Live Agent: Human-in-the-Loop</Heading>
            <p>
              Deploy persistent AI agents your team can interact with in real time
              &mdash; through the web terminal, CLI, or by submitting Tasks.
            </p>
            <ul>
              <li>Zero cold start &mdash; agent is always running</li>
              <li>Interactive terminal access via CLI or web</li>
              <li>Auto-suspend when idle, resume on demand</li>
              <li>Session history persists across restarts</li>
            </ul>
          </div>
          <div className="col col--6">
            <CodeBlock language="yaml" title="agent.yaml">
              {agentYaml}
            </CodeBlock>
          </div>
        </div>
        <div className="row" style={{marginTop: '2rem'}}>
          <div className="col col--6">
            <Heading as="h2">AgentTemplate: Workflows at Scale</Heading>
            <p>
              Run stable, repeatable AI tasks in ephemeral Pods. Perfect for
              CI/CD pipelines, batch operations, and automated workflows.
            </p>
            <ul>
              <li>No new tools to learn &mdash; just <code>kubectl apply</code></li>
              <li>Works with any CI/CD pipeline</li>
              <li>Scale with Helm templates for batch operations</li>
              <li>Rate limiting and quota controls</li>
            </ul>
          </div>
          <div className="col col--6">
            <CodeBlock language="yaml" title="task.yaml">
              {taskYaml}
            </CodeBlock>
          </div>
        </div>
      </div>
    </section>
  );
}

function StarCallToAction() {
  return (
    <section className={styles.starCta}>
      <div className="container text--center">
        <Heading as="h2">Like this project?</Heading>
        <p>
          If you find KubeOpenCode useful, please give us a star on GitHub.
          It helps others discover the project and motivates us to keep improving it.
        </p>
        <Link
          className="button button--primary button--lg"
          href={`https://github.com/${GITHUB_REPO}`}>
          Star on GitHub
        </Link>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="Kubernetes-native Agent Platform for Teams and Enterprise"
      description="Deploy, manage, and govern AI coding agents at scale on Kubernetes. Built on OpenCode, designed for teams and enterprise.">
      <AlphaBanner />
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <QuickExample />
        <StarCallToAction />
      </main>
    </Layout>
  );
}
