import { describe, it, expect, vi } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithProviders } from '../../test/utils';
import YamlViewer from '../YamlViewer';

// Mock CodeMirror since it requires a real DOM
vi.mock('@uiw/react-codemirror', () => ({
  __esModule: true,
  default: ({ value, onChange, extensions }: { value: string; onChange?: (v: string) => void; extensions?: unknown[] }) => {
    const isReadOnly = extensions?.some((ext: unknown) => {
      // Check if EditorView.editable.of(false) is in extensions
      return ext && typeof ext === 'object' && 'value' in (ext as Record<string, unknown>);
    });
    return isReadOnly ? (
      <pre data-testid="codemirror-readonly">{value}</pre>
    ) : (
      <textarea
        data-testid="codemirror-editor"
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
      />
    );
  },
}));

vi.mock('@codemirror/lang-yaml', () => ({
  yaml: () => ({}),
}));

vi.mock('@codemirror/view', () => ({
  EditorView: {
    theme: () => ({}),
    editable: { of: (v: boolean) => ({ value: v }) },
  },
}));

const sampleYaml = `apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: test-task
  namespace: default`;

describe('YamlViewer', () => {
  it('renders collapsed by default', () => {
    renderWithProviders(
      <YamlViewer queryKey={['test']} fetchYaml={() => Promise.resolve(sampleYaml)} />
    );

    expect(screen.getByText('YAML')).toBeInTheDocument();
    expect(screen.queryByText('Resource Definition')).not.toBeInTheDocument();
  });

  it('expands when YAML button is clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer queryKey={['test']} fetchYaml={() => Promise.resolve(sampleYaml)} />
    );

    await user.click(screen.getByText('YAML'));

    await waitFor(() => {
      expect(screen.getByText('Resource Definition')).toBeInTheDocument();
    });
  });

  it('shows YAML content after expanding', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer queryKey={['test']} fetchYaml={() => Promise.resolve(sampleYaml)} />
    );

    await user.click(screen.getByText('YAML'));

    await waitFor(() => {
      const pre = screen.getByTestId('codemirror-readonly');
      expect(pre).toBeInTheDocument();
      expect(pre.textContent).toContain('apiVersion: kubeopencode.io/v1alpha1');
      expect(pre.textContent).toContain('kind: Task');
    });
  });

  it('shows Copy button when YAML is loaded', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer queryKey={['test']} fetchYaml={() => Promise.resolve(sampleYaml)} />
    );

    await user.click(screen.getByText('YAML'));

    await waitFor(() => {
      expect(screen.getByText('Copy')).toBeInTheDocument();
    });
  });

  it('shows loading state while fetching', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer
        queryKey={['test']}
        fetchYaml={() => new Promise((resolve) => setTimeout(() => resolve(sampleYaml), 100))}
      />
    );

    await user.click(screen.getByText('YAML'));

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('shows error state when fetch fails', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer
        queryKey={['test-error']}
        fetchYaml={() => Promise.reject(new Error('Network error'))}
      />
    );

    await user.click(screen.getByText('YAML'));

    await waitFor(() => {
      expect(screen.getByText(/Error: Network error/)).toBeInTheDocument();
    });
  });

  it('collapses when YAML button is clicked again', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer queryKey={['test']} fetchYaml={() => Promise.resolve(sampleYaml)} />
    );

    // Expand
    await user.click(screen.getByText('YAML'));
    await waitFor(() => {
      expect(screen.getByText('Resource Definition')).toBeInTheDocument();
    });

    // Collapse
    await user.click(screen.getByText('YAML'));
    expect(screen.queryByText('Resource Definition')).not.toBeInTheDocument();
  });

  it('shows Edit button when onSave is provided', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer
        queryKey={['test']}
        fetchYaml={() => Promise.resolve(sampleYaml)}
        onSave={async () => {}}
      />
    );

    await user.click(screen.getByText('YAML'));

    await waitFor(() => {
      expect(screen.getByText('Edit')).toBeInTheDocument();
    });
  });

  it('shows Format button in edit mode', async () => {
    const user = userEvent.setup();
    renderWithProviders(
      <YamlViewer
        queryKey={['test']}
        fetchYaml={() => Promise.resolve(sampleYaml)}
        onSave={async () => {}}
      />
    );

    await user.click(screen.getByText('YAML'));
    await waitFor(() => {
      expect(screen.getByText('Edit')).toBeInTheDocument();
    });

    await user.click(screen.getByText('Edit'));

    expect(screen.getByText('Format')).toBeInTheDocument();
    expect(screen.getByText('Save')).toBeInTheDocument();
    expect(screen.getByText('Cancel')).toBeInTheDocument();
  });
});
