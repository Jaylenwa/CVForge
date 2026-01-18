import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Search } from 'lucide-react';

export type DataTableFixed = 'right';

export interface DataTableColumn<T> {
  key: string;
  label: React.ReactNode;
  minWidth?: string | number;
  fixed?: DataTableFixed;
  headerClassName?: string;
  cellClassName?: string;
  nowrap?: boolean;
  render?: (row: T) => React.ReactNode;
}

export interface DataTableEmptyState {
  title: string;
  description?: string;
}

export interface DataTableProps<T> {
  data: T[];
  columns: Array<DataTableColumn<T>>;
  getRowKey: (row: T, index: number) => string;
  emptyState?: DataTableEmptyState;
  rowClassName?: (row: T, index: number) => string;
}

export const DataTable = <T,>({ data, columns, getRowKey, emptyState, rowClassName }: DataTableProps<T>) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isScrolledRight, setIsScrolledRight] = useState(false);

  const fixedRightColumns = useMemo(() => columns.filter(c => c.fixed === 'right'), [columns]);
  const hasFixedRight = fixedRightColumns.length > 0;

  const handleScroll = () => {
    const el = containerRef.current;
    if (!el || !hasFixedRight) return;
    const { scrollLeft, scrollWidth, clientWidth } = el;
    const atEnd = scrollLeft + clientWidth >= scrollWidth - 10;
    setIsScrolledRight(!atEnd);
  };

  useEffect(() => {
    handleScroll();
    window.addEventListener('resize', handleScroll);
    return () => window.removeEventListener('resize', handleScroll);
  }, [data.length, hasFixedRight]);

  const rightOffsets = useMemo(() => {
    const offsets: Record<string, number> = {};
    let acc = 0;
    for (let i = columns.length - 1; i >= 0; i--) {
      const col = columns[i];
      if (col.fixed !== 'right') continue;
      offsets[col.key] = acc;
      const w = col.minWidth;
      if (typeof w === 'number') acc += w;
      else if (typeof w === 'string' && w.endsWith('px')) acc += Number(w.replace('px', '')) || 0;
      else acc += 140;
    }
    return offsets;
  }, [columns]);

  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      className="relative flex-1 overflow-x-auto"
    >
      <table className="w-full border-collapse text-left text-sm">
        <thead className="sticky top-0 z-20 bg-gray-50 shadow-sm">
          <tr>
            {columns.map((col) => {
              const isFixedRight = col.fixed === 'right';
              const right = isFixedRight ? rightOffsets[col.key] || 0 : undefined;
              return (
                <th
                  key={col.key}
                  style={{ minWidth: col.minWidth, right }}
                  className={[
                    'px-6 py-4 font-semibold text-gray-600 border-b border-gray-100',
                    isFixedRight ? 'sticky bg-gray-50 z-30' : '',
                    col.headerClassName || '',
                  ].join(' ')}
                >
                  <div className="flex items-center gap-2">
                    {col.label}
                    {isFixedRight && isScrolledRight ? (
                      <div className="absolute left-0 top-0 bottom-0 w-8 -translate-x-full pointer-events-none bg-gradient-to-r from-transparent to-black/5" />
                    ) : null}
                  </div>
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-50">
          {data.map((row, idx) => (
            <tr key={getRowKey(row, idx)} className={['hover:bg-blue-50/30 transition-colors group', rowClassName?.(row, idx) || ''].join(' ')}>
              {columns.map((col) => {
                const isFixedRight = col.fixed === 'right';
                const right = isFixedRight ? rightOffsets[col.key] || 0 : undefined;
                return (
                  <td
                    key={col.key}
                    style={{ right }}
                    className={[
                      'px-6 py-3 align-top',
                      col.nowrap ? 'whitespace-nowrap' : '',
                      isFixedRight ? 'sticky bg-white group-hover:bg-blue-50/50 z-10' : '',
                      col.cellClassName || '',
                    ].join(' ')}
                  >
                    {isFixedRight && isScrolledRight ? (
                      <div className="absolute left-0 top-0 bottom-0 w-8 -translate-x-full pointer-events-none bg-gradient-to-r from-transparent to-black/5" />
                    ) : null}
                    {col.render ? col.render(row) : null}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>

      {data.length === 0 && emptyState ? (
        <div className="flex flex-col items-center justify-center py-20 text-gray-500">
          <div className="p-4 bg-gray-50 rounded-full mb-4">
            <Search size={32} className="text-gray-300" />
          </div>
          <p className="text-lg font-medium">{emptyState.title}</p>
          {emptyState.description ? <p className="text-sm">{emptyState.description}</p> : null}
        </div>
      ) : null}
    </div>
  );
};

