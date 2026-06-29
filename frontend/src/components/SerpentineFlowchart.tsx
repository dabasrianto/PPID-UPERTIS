import React from 'react';
import { ArrowRight, ArrowLeft, ArrowDown } from 'lucide-react';

export interface FlowchartStep {
  step: number;
  title: string;
  desc: string;
  icon: React.ComponentType<any>;
  color: string;
  bullets?: string[];
}

interface SerpentineFlowchartProps {
  steps: FlowchartStep[];
  activeStepIdx: number;
  setActiveStepIdx: (idx: number) => void;
}

export default function SerpentineFlowchart({
  steps,
  activeStepIdx,
  setActiveStepIdx
}: SerpentineFlowchartProps) {
  if (!steps || steps.length === 0) return null;

  const chunkArray = (arr: any[], size: number) => {
    const chunks = [];
    for (let i = 0; i < arr.length; i += size) {
      chunks.push(arr.slice(i, i + size));
    }
    return chunks;
  };

  const stepRows = chunkArray(steps, 4);

  return (
    <div className="w-full">
      {/* Desktop Layout: Horizontal Pipeline Wrapping Grid */}
      <div className="hidden md:block space-y-6">
        {stepRows.map((row, rowIdx) => {
          const isEven = rowIdx % 2 === 0;
          const displayItems = isEven ? row : [...row].reverse();
          const spacerCount = 4 - row.length;
          
          return (
            <div key={rowIdx} className="space-y-6">
              {/* Downward connector between rows */}
              {rowIdx > 0 && (
                <div className={`flex ${rowIdx % 2 === 1 ? 'justify-end pr-14' : 'justify-start pl-14'} -my-2`}>
                  <div className="bg-white p-1 rounded-full border border-slate-150 shadow-sm text-slate-400">
                    <ArrowDown className="h-3.5 w-3.5" />
                  </div>
                </div>
              )}

              {/* Grid row */}
              <div className="grid grid-cols-4 gap-4 items-center relative">
                {/* If odd row (reversed), prepend empty spacers for correct right-alignment under preceding row */}
                {!isEven && Array.from({ length: spacerCount }).map((_, sIdx) => (
                  <div key={`spacer-${sIdx}`} className="col-span-1" />
                ))}

                {displayItems.map((item, colIdx) => {
                  const IconComponent = item.icon;
                  const originalIdx = steps.findIndex(s => s.step === item.step);
                  const isActive = activeStepIdx === originalIdx;
                  
                  // In reversed row (odd), index flows right to left.
                  // For normal row (even), index flows left to right.
                  const isLastInGridRow = colIdx === displayItems.length - 1;
                  
                  return (
                    <div key={item.step} className="flex items-center w-full relative">
                      <button
                        onClick={() => setActiveStepIdx(originalIdx)}
                        className={`flex-1 flex flex-col items-center justify-center border rounded-2xl p-4 text-center transition-all cursor-pointer h-36 relative ${
                          isActive
                            ? 'border-amber-400 bg-amber-500/5 ring-2 ring-amber-400/20'
                            : 'border-slate-200 hover:border-amber-400 hover:bg-slate-50/50'
                        }`}
                      >
                        <span className={`text-[10px] font-extrabold px-2 py-0.5 rounded-full absolute top-2 left-2 ${isActive ? 'bg-amber-400 text-[#002147]' : 'bg-slate-100 text-slate-400'}`}>
                          Langkah {item.step}
                        </span>
                        <div className={`h-11 w-11 rounded-xl flex items-center justify-center border-2 ${item.color} bg-white shadow-sm mb-2`}>
                          <IconComponent className="h-5 w-5" />
                        </div>
                        <span className="text-[11px] font-bold text-slate-800 leading-snug truncate w-full px-1">{item.title}</span>
                      </button>

                      {/* Right pointing connector for even rows */}
                      {isEven && !isLastInGridRow && (
                        <div className="absolute -right-3 top-1/2 -translate-y-1/2 z-20 bg-white p-0.5 rounded-full border border-slate-100 shadow-sm text-slate-400">
                          <ArrowRight className="h-3.5 w-3.5" />
                        </div>
                      )}

                      {/* Left pointing connector for odd rows */}
                      {!isEven && !isLastInGridRow && (
                        <div className="absolute -left-3 top-1/2 -translate-y-1/2 z-20 bg-white p-0.5 rounded-full border border-slate-100 shadow-sm text-slate-400">
                          <ArrowLeft className="h-3.5 w-3.5" />
                        </div>
                      )}
                    </div>
                  );
                })}

                {/* If even row (normal), append empty spacers for padding if it's the last row */}
                {isEven && Array.from({ length: spacerCount }).map((_, sIdx) => (
                  <div key={`spacer-even-${sIdx}`} className="col-span-1" />
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* Mobile Layout: Clean Vertical Pipeline */}
      <div className="block md:hidden space-y-4">
        {steps.map((item: any, idx: number) => {
          const IconComponent = item.icon;
          const isActive = activeStepIdx === idx;
          return (
            <div key={idx} className="flex flex-col items-center w-full">
              <button
                onClick={() => setActiveStepIdx(idx)}
                className={`w-full flex items-center gap-4 border rounded-2xl p-4 transition-all cursor-pointer relative ${
                  isActive
                    ? 'border-amber-400 bg-amber-500/5 ring-2 ring-amber-400/20'
                    : 'border-slate-200 bg-white'
                }`}
              >
                <div className={`h-10 w-10 rounded-lg flex items-center justify-center border-2 ${item.color} bg-white shadow-sm shrink-0`}>
                  <IconComponent className="h-5 w-5" />
                </div>
                <div className="text-left">
                  <span className="text-[9px] font-bold text-slate-400 uppercase tracking-widest block">Langkah {item.step}</span>
                  <span className="text-xs font-bold text-slate-800">{item.title}</span>
                </div>
              </button>
              {idx < steps.length - 1 && (
                <div className="py-1 text-slate-350">
                  <ArrowDown className="h-4 w-4" />
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
