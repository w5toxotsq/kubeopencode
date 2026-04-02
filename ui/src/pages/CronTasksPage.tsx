import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../api/client';
import TimeAgo from '../components/TimeAgo';
import ResourceFilter from '../components/ResourceFilter';
import { TableSkeleton } from '../components/Skeleton';
import ConfirmDialog from '../components/ConfirmDialog';
import { useFilterState } from '../hooks/useFilterState';
import { useNamespace } from '../contexts/NamespaceContext';
import { useToast } from '../contexts/ToastContext';

const PAGE_SIZE_OPTIONS = [10, 20, 50];

function CronTasksPage() {
  const { namespace, isAllNamespaces } = useNamespace();
  const queryClient = useQueryClient();
  const { addToast } = useToast();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [filters, setFilters] = useFilterState();
  const [deleteTarget, setDeleteTarget] = useState<{ namespace: string; name: string } | null>(null);

  useEffect(() => {
    setCurrentPage(1);
  }, [namespace, filters.name, filters.labelSelector]);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['crontasks', namespace, currentPage, pageSize, filters.name, filters.labelSelector],
    queryFn: () => {
      const params = {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
        sortOrder: 'desc' as const,
        name: filters.name || undefined,
        labelSelector: filters.labelSelector || undefined,
      };
      return isAllNamespaces
        ? api.listAllCronTasks(params)
        : api.listCronTasks(namespace, params);
    },
    refetchInterval: 5000,
  });

  const deleteMutation = useMutation({
    mutationFn: ({ ns, n }: { ns: string; n: string }) => api.deleteCronTask(ns, n),
    onSuccess: (_, { n }) => {
      addToast(`CronTask "${n}" deleted successfully`, 'success');
      queryClient.invalidateQueries({ queryKey: ['crontasks'] });
    },
    onError: (err: Error) => {
      addToast(`Failed to delete CronTask: ${err.message}`, 'error');
    },
  });

  const triggerMutation = useMutation({
    mutationFn: ({ ns, n }: { ns: string; n: string }) => api.triggerCronTask(ns, n),
    onSuccess: (_, { n }) => {
      addToast(`CronTask "${n}" triggered successfully`, 'success');
      queryClient.invalidateQueries({ queryKey: ['crontasks'] });
    },
    onError: (err: Error) => {
      addToast(`Failed to trigger CronTask: ${err.message}`, 'error');
    },
  });

  const suspendMutation = useMutation({
    mutationFn: ({ ns, n }: { ns: string; n: string }) => api.suspendCronTask(ns, n),
    onSuccess: (_, { n }) => {
      addToast(`CronTask "${n}" suspended`, 'success');
      queryClient.invalidateQueries({ queryKey: ['crontasks'] });
    },
    onError: (err: Error) => {
      addToast(`Failed to suspend CronTask: ${err.message}`, 'error');
    },
  });

  const resumeMutation = useMutation({
    mutationFn: ({ ns, n }: { ns: string; n: string }) => api.resumeCronTask(ns, n),
    onSuccess: (_, { n }) => {
      addToast(`CronTask "${n}" resumed`, 'success');
      queryClient.invalidateQueries({ queryKey: ['crontasks'] });
    },
    onError: (err: Error) => {
      addToast(`Failed to resume CronTask: ${err.message}`, 'error');
    },
  });

  return (
    <div className="animate-fade-in">
      <div className="sm:flex sm:items-center sm:justify-between mb-6">
        <div>
          <h2 className="font-display text-2xl font-bold text-stone-900 tracking-tight">CronTasks</h2>
          <p className="mt-1 text-sm text-stone-500">
            Scheduled AI agent tasks running on a cron schedule
          </p>
        </div>
        <div className="mt-4 sm:mt-0">
          <Link
            to="/crontasks/create"
            className="inline-flex items-center gap-2 px-4 py-2.5 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors shadow-sm"
          >
            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 5v14M5 12h14" strokeLinecap="round" />
            </svg>
            New CronTask
          </Link>
        </div>
      </div>

      {/* Filter bar */}
      <div className="mb-4">
        <ResourceFilter
          filters={filters}
          onFilterChange={setFilters}
          placeholder="Filter CronTasks by name..."
        />
      </div>

      {isLoading ? (
        <div className="bg-white rounded-xl border border-stone-200 overflow-hidden shadow-sm">
          <TableSkeleton rows={5} cols={isAllNamespaces ? 9 : 8} />
        </div>
      ) : error ? (
        <div className="bg-red-50 border border-red-200 rounded-xl p-5">
          <p className="text-red-700 text-sm">Error loading CronTasks: {(error as Error).message}</p>
          <button
            onClick={() => refetch()}
            className="mt-2 text-sm text-red-600 hover:text-red-800 font-medium"
          >
            Retry
          </button>
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-stone-200 overflow-hidden shadow-sm">
          <table className="min-w-full divide-y divide-stone-100">
            <thead className="bg-stone-50/60">
              <tr>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Name
                </th>
                {isAllNamespaces && (
                  <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                    Namespace
                  </th>
                )}
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Schedule
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Source
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider hidden sm:table-cell">
                  Active
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider hidden lg:table-cell">
                  Last Run
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider hidden lg:table-cell">
                  Next Run
                </th>
                <th className="px-5 py-3 text-left text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Age
                </th>
                <th className="px-5 py-3 text-right text-[11px] font-display font-medium text-stone-400 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-stone-100">
              {data?.cronTasks.length === 0 ? (
                <tr>
                  <td colSpan={isAllNamespaces ? 10 : 9} className="px-5 py-12 text-center text-stone-400 text-sm">
                    No CronTasks found.{' '}
                    {!isAllNamespaces && (
                      <Link to="/crontasks/create" className="text-primary-600 hover:text-primary-700 font-medium">
                        Create your first CronTask
                      </Link>
                    )}
                  </td>
                </tr>
              ) : (
                data?.cronTasks.map((ct) => (
                  <tr key={`${ct.namespace}/${ct.name}`} className="hover:bg-stone-50/60 transition-colors">
                    <td className="px-5 py-3.5 whitespace-nowrap">
                      <Link
                        to={`/crontasks/${ct.namespace}/${ct.name}`}
                        className="text-stone-800 hover:text-primary-600 font-medium text-sm transition-colors"
                      >
                        {ct.name}
                      </Link>
                    </td>
                    {isAllNamespaces && (
                      <td className="px-5 py-3.5 whitespace-nowrap text-sm text-stone-400">
                        {ct.namespace}
                      </td>
                    )}
                    <td className="px-5 py-3.5 whitespace-nowrap text-xs font-mono text-stone-600">
                      {ct.schedule}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-xs">
                      {ct.taskTemplate.agentRef ? (
                        <Link to={`/agents/${ct.namespace}/${ct.taskTemplate.agentRef.name}`} className="text-stone-500 hover:text-primary-600 font-mono transition-colors">
                          {ct.taskTemplate.agentRef.name}
                        </Link>
                      ) : ct.taskTemplate.templateRef ? (
                        <Link to={`/templates/${ct.namespace}/${ct.taskTemplate.templateRef.name}`} className="text-amber-600 hover:text-amber-700 font-mono transition-colors">
                          {ct.taskTemplate.templateRef.name}
                        </Link>
                      ) : (
                        <span className="text-stone-400 font-mono">-</span>
                      )}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap">
                      {ct.suspend ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[11px] font-medium bg-stone-100 text-stone-600 border border-stone-200">
                          Suspended
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[11px] font-medium bg-emerald-50 text-emerald-700 border border-emerald-200">
                          Active
                        </span>
                      )}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-sm text-stone-600 hidden sm:table-cell font-mono text-xs">
                      {ct.active}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-xs text-stone-400 hidden lg:table-cell">
                      {ct.lastScheduleTime ? <TimeAgo date={ct.lastScheduleTime} /> : '-'}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-xs text-stone-400 hidden lg:table-cell">
                      {ct.nextScheduleTime ? <TimeAgo date={ct.nextScheduleTime} /> : '-'}
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-xs text-stone-400">
                      <TimeAgo date={ct.createdAt} />
                    </td>
                    <td className="px-5 py-3.5 whitespace-nowrap text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <button
                          onClick={() => triggerMutation.mutate({ ns: ct.namespace, n: ct.name })}
                          disabled={triggerMutation.isPending}
                          className="px-2.5 py-1 text-[11px] font-medium text-primary-700 bg-primary-50 border border-primary-200 rounded-md hover:bg-primary-100 transition-colors"
                          title="Run Now"
                        >
                          Run
                        </button>
                        {ct.suspend ? (
                          <button
                            onClick={() => resumeMutation.mutate({ ns: ct.namespace, n: ct.name })}
                            disabled={resumeMutation.isPending}
                            className="px-2.5 py-1 text-[11px] font-medium text-emerald-700 bg-emerald-50 border border-emerald-200 rounded-md hover:bg-emerald-100 transition-colors"
                            title="Resume"
                          >
                            Resume
                          </button>
                        ) : (
                          <button
                            onClick={() => suspendMutation.mutate({ ns: ct.namespace, n: ct.name })}
                            disabled={suspendMutation.isPending}
                            className="px-2.5 py-1 text-[11px] font-medium text-amber-700 bg-amber-50 border border-amber-200 rounded-md hover:bg-amber-100 transition-colors"
                            title="Suspend"
                          >
                            Suspend
                          </button>
                        )}
                        <button
                          onClick={() => setDeleteTarget({ namespace: ct.namespace, name: ct.name })}
                          className="px-2.5 py-1 text-[11px] font-medium text-red-600 bg-red-50 border border-red-200 rounded-md hover:bg-red-100 transition-colors"
                          title="Delete"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>

          {/* Pagination */}
          {data?.pagination && data.pagination.totalCount > 0 && (
            <div className="bg-white px-5 py-3 flex items-center justify-between border-t border-stone-100">
              <div className="flex-1 flex justify-between sm:hidden">
                <button
                  onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className="px-3 py-1.5 text-sm font-medium text-stone-600 bg-stone-100 rounded-lg hover:bg-stone-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  Previous
                </button>
                <button
                  onClick={() => setCurrentPage(p => p + 1)}
                  disabled={!data.pagination.hasMore}
                  className="ml-3 px-3 py-1.5 text-sm font-medium text-stone-600 bg-stone-100 rounded-lg hover:bg-stone-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  Next
                </button>
              </div>
              <div className="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                <div className="flex items-center space-x-4">
                  <p className="text-xs text-stone-400">
                    <span className="font-medium text-stone-600">{data.pagination.offset + 1}</span>
                    {' '}-{' '}
                    <span className="font-medium text-stone-600">
                      {Math.min(data.pagination.offset + data.cronTasks.length, data.pagination.totalCount)}
                    </span>
                    {' '}of{' '}
                    <span className="font-medium text-stone-600">{data.pagination.totalCount}</span>
                  </p>
                  <select
                    value={pageSize}
                    onChange={(e) => {
                      setPageSize(Number(e.target.value));
                      setCurrentPage(1);
                    }}
                    className="block w-16 rounded-lg border-stone-200 text-xs text-stone-600 focus:border-primary-500 focus:ring-primary-500"
                  >
                    {PAGE_SIZE_OPTIONS.map((size) => (
                      <option key={size} value={size}>{size}</option>
                    ))}
                  </select>
                </div>
                <div className="flex items-center space-x-1">
                  <button
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="px-3 py-1.5 text-xs font-medium text-stone-500 bg-stone-50 border border-stone-200 rounded-lg hover:bg-stone-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    Prev
                  </button>
                  <span className="px-3 py-1.5 text-xs font-mono text-stone-500">
                    {currentPage}/{Math.ceil(data.pagination.totalCount / pageSize)}
                  </span>
                  <button
                    onClick={() => setCurrentPage(p => p + 1)}
                    disabled={!data.pagination.hasMore}
                    className="px-3 py-1.5 text-xs font-medium text-stone-500 bg-stone-50 border border-stone-200 rounded-lg hover:bg-stone-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        title="Delete CronTask"
        message={`Are you sure you want to delete CronTask "${deleteTarget?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
        variant="danger"
        onConfirm={() => {
          if (deleteTarget) {
            deleteMutation.mutate({ ns: deleteTarget.namespace, n: deleteTarget.name });
          }
          setDeleteTarget(null);
        }}
        onCancel={() => setDeleteTarget(null)}
      />
    </div>
  );
}

export default CronTasksPage;
