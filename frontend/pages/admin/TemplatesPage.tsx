import React, { useState } from 'react';
import { createTemplate, updateTemplate, deleteTemplate } from '../../services/adminService';
import { Button } from '../../components/ui/Button';
import { useToast } from '../../components/ui/Toast';
import { useConfirm } from '../../components/ui/ConfirmDialog';
import { useLanguage } from '../../contexts/LanguageContext';

export const TemplatesPage: React.FC = () => {
  const [form, setForm] = useState({ externalId: '', name: '', thumbnail: '', tags: '', category: '', level: '', popularity: 0, isPremium: false });
  const { showToast } = useToast();
  const confirm = useConfirm();
  const { t } = useLanguage();

  const create = async () => {
    try {
      await createTemplate({ ...form });
      showToast(t('admin.msg.templateCreated'), 'success');
    } catch {
      showToast(t('admin.msg.templateCreateFailed'), 'error');
    }
  };

  const patch = async () => {
    if (!form.externalId) return;
    try {
      await updateTemplate(form.externalId, { name: form.name, thumbnail: form.thumbnail, tags: form.tags, category: form.category, level: form.level, popularity: form.popularity, isPremium: form.isPremium });
      showToast(t('admin.msg.templateUpdated'), 'success');
    } catch {
      showToast(t('admin.msg.templateUpdateFailed'), 'error');
    }
  };

  const remove = async () => {
    if (!form.externalId) return;
    const ok = await confirm({ title: t('common.confirmAction'), message: t('admin.actions.delete') });
    if (!ok) return;
    try {
      await deleteTemplate(form.externalId);
      showToast(t('admin.msg.templateDeleted'), 'success');
    } catch {
      showToast(t('admin.msg.templateDeleteFailed'), 'error');
    }
  };

  return (
    <div className="space-y-4">
      <div className="bg-white shadow-sm rounded-md p-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.externalId')} value={form.externalId} onChange={e => setForm({ ...form, externalId: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.name')} value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.thumbnail')} value={form.thumbnail} onChange={e => setForm({ ...form, thumbnail: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.tags')} value={form.tags} onChange={e => setForm({ ...form, tags: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.category')} value={form.category} onChange={e => setForm({ ...form, category: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" placeholder={t('admin.form.level')} value={form.level} onChange={e => setForm({ ...form, level: e.target.value })} />
          <input className="border rounded-md px-3 py-2 text-sm" type="number" placeholder={t('admin.form.popularity')} value={form.popularity} onChange={e => setForm({ ...form, popularity: parseInt(e.target.value || '0') })} />
          <label className="inline-flex items-center space-x-2 text-sm">
            <input type="checkbox" checked={form.isPremium} onChange={e => setForm({ ...form, isPremium: e.target.checked })} />
            <span>{t('admin.form.isPremium')}</span>
          </label>
        </div>
        <div className="mt-4 space-x-2">
          <Button onClick={create}>{t('admin.actions.create')}</Button>
          <Button variant="outline" onClick={patch}>{t('admin.actions.update')}</Button>
          <Button variant="danger" onClick={remove}>{t('admin.actions.delete')}</Button>
        </div>
      </div>
    </div>
  );
};
