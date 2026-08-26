export function paraDataApi(dataInput: string): string {
  return `${dataInput}T00:00:00Z`;
}

export function paraDataInput(dataIso: string): string {
  return dataIso.slice(0, 10);
}

export function formatarDataLocal(data: Date): string {
  const ano = data.getFullYear();
  const mes = String(data.getMonth() + 1).padStart(2, '0');
  const dia = String(data.getDate()).padStart(2, '0');
  return `${ano}-${mes}-${dia}`;
}
