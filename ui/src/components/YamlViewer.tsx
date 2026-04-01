import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';

interface YamlViewerProps {
  queryKey: string[];
  fetchYaml: () => Promise<string>;
  onSave?: (yaml: string) => Promise<void>;
}

function YamlViewer({ queryKey, fetchYaml, onSave }: YamlViewerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [editing, setEditing] = useState(false);
  const [editValue, setEditValue] = useState('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  const { data: yaml, isLoading, error, refetch } = useQuery({
    queryKey: [...queryKey, 'yaml'],
    queryFn: fetchYaml,
    enabled: isOpen,
  });

  useEffect(() => {
    if (yaml && editing) {
      setEditValue(yaml);
    }
  }, [yaml, editing]);

  const handleEdit = () => {
    setEditValue(yaml || '');
    setEditing(true);
    setSaveError('');
  };

  const handleCancel = () => {
    setEditing(false);
    setSaveError('');
  };

  const handleSave = async () => {
    if (!onSave) return;
    setSaving(true);
    setSaveError('');
    try {
      await onSave(editValue);
      setEditing(false);
      refetch();
    } catch (err) {
      setSaveError((err as Error).message || 'Failed to save');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="mt-6">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 text-xs font-display font-medium text-stone-400 hover:text-stone-600 uppercase tracking-wider transition-colors"
      >
        <svg
          className={`w-3.5 h-3.5 transform transition-transform ${isOpen ? 'rotate-90' : ''}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
        <span>YAML</span>
      </button>
      {isOpen && (
        <div className="mt-2 bg-stone-900 rounded-xl overflow-hidden border border-stone-800 animate-fade-in">
          <div className="px-4 py-2.5 bg-stone-800/50 flex items-center justify-between border-b border-stone-700/50">
            <span className="text-xs text-stone-400 font-display">Resource Definition</span>
            <div className="flex items-center gap-2">
              {yaml && !editing && (
                <button
                  onClick={() => navigator.clipboard.writeText(yaml)}
                  className="text-[11px] text-stone-500 hover:text-stone-300 transition-colors font-medium"
                >
                  Copy
                </button>
              )}
              {onSave && !editing && yaml && (
                <button
                  onClick={handleEdit}
                  className="text-[11px] text-sky-400 hover:text-sky-300 transition-colors font-medium"
                >
                  Edit
                </button>
              )}
            </div>
          </div>
          <div className="p-4">
            {isLoading ? (
              <span className="text-stone-500 text-sm">Loading...</span>
            ) : error ? (
              <span className="text-red-400 text-sm">Error: {(error as Error).message}</span>
            ) : editing ? (
              <div>
                <textarea
                  value={editValue}
                  onChange={(e) => setEditValue(e.target.value)}
                  className="w-full h-96 bg-stone-950 text-stone-300 font-mono text-xs p-3 rounded-lg border border-stone-700 focus:border-sky-500 focus:outline-none resize-y leading-relaxed"
                  spellCheck={false}
                />
                {saveError && (
                  <p className="text-xs text-red-400 mt-2">{saveError}</p>
                )}
                <div className="flex items-center justify-end gap-2 mt-3">
                  <button
                    onClick={handleCancel}
                    disabled={saving}
                    className="px-3 py-1.5 text-xs font-medium text-stone-400 bg-stone-800 rounded-lg hover:bg-stone-700 transition-colors disabled:opacity-50"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleSave}
                    disabled={saving}
                    className="px-3 py-1.5 text-xs font-medium text-white bg-sky-600 rounded-lg hover:bg-sky-700 transition-colors disabled:opacity-50"
                  >
                    {saving ? 'Saving...' : 'Save'}
                  </button>
                </div>
              </div>
            ) : (
              <div className="max-h-96 overflow-y-auto sidebar-scroll">
                <pre className="text-xs text-stone-300 font-mono whitespace-pre leading-relaxed">{yaml}</pre>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default YamlViewer;
