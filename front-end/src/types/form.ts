export type Field = {
  id: string;
  label: string;
  placeholder: string;
  variant: 'primary';
  type?: string;
  className?: string;
};

export type SubmitButton = {
    content: string;
    variant: 'brand' | 'neutral';
}