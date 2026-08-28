/** Buduje wiersz rejestru akcji, który nie jest operacją kontekstową i nie ma wybieraka, bo identyfikator wskazuje komendę okna komunikacji. */
export function utworzPozycjeSpozaOperacji(
  id: string,
  nazwa: string,
  komenda: string,
): HTMLElement {
  const etykieta = document.createElement('span');
  etykieta.className = 'ms-wykaz__etykieta';
  etykieta.textContent = nazwa;

  const zrodlo = document.createElement('span');
  zrodlo.className = 'dn-plakietka ms-wykaz__zrodlo';
  zrodlo.textContent = `rejestr rdzenia · komenda ${komenda}`;
  zrodlo.title = `Wiersz rejestru ${id} wskazuje komendę ${komenda}, nie operację kontekstową Studio.`;

  const wiersz = document.createElement('div');
  wiersz.className = 'ms-wykaz__wiersz';
  wiersz.dataset['operacja'] = id;
  wiersz.dataset['wybieralna'] = 'nie';
  wiersz.append(etykieta, zrodlo);
  return wiersz;
}
