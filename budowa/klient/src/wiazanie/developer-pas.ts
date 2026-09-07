// Wspólny pas działań okna Developer: pole wpisu, przyciski czynności,
// uzbrajanie czynności nieodwracalnych i wypełnianie ciała panelu.

/*
zalozPas stawia pas nad ciałem wskazanego panelu.

Pas stoi nad ciałem, bo ciało jest wymieniane przy każdym odświeżeniu wykazu;
pas postawiony w nim znikałby razem z pozycjami.
*/
export function zalozPas(
  korzen: Element,
  panelKod: string,
  znacznik: string,
  podpowiedz: string,
  przyciski: ReadonlyArray<readonly [string, string]>,
): void {
  const panel = korzen.querySelector(`#${panelKod}`);
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  if (panel.querySelector(`[data-${znacznik}-wpis]`) !== null) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = korzen.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = podpowiedz;
  pole.setAttribute('aria-label', podpowiedz);
  pole.setAttribute(`data-${znacznik}-wpis`, '');
  pas.appendChild(pole);
  for (const [czynnosc, etykieta] of przyciski) {
    const wezel = korzen.ownerDocument.createElement('button');
    wezel.type = 'button';
    wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
    wezel.setAttribute(`data-${znacznik}`, czynnosc);
    wezel.textContent = etykieta;
    pas.appendChild(wezel);
  }
  panel.insertBefore(pas, cialo);
}

export function wpisPasa(korzen: Element, znacznik: string): string {
  const pole = korzen.querySelector<HTMLInputElement>(`[data-${znacznik}-wpis]`);
  return (pole?.value ?? '').trim();
}

export function czesciWpisu(korzen: Element, znacznik: string): string[] {
  return wpisPasa(korzen, znacznik).split('|').map((czesc) => czesc.trim());
}

/* Drugie naciśnięcie potwierdza czynność nieodwracalną; napis wraca sam, więc
   uzbrojenie nie zostaje na przycisku na stałe. */
export function potwierdzone(korzen: Element, znacznik: string, czynnosc: string): boolean {
  const guzik = korzen.querySelector<HTMLElement>(`[data-${znacznik}="${czynnosc}"]`);
  if (guzik === null) return true;
  if (guzik.dataset.uzbrojone === 'tak') {
    delete guzik.dataset.uzbrojone;
    return true;
  }
  const napis = guzik.textContent ?? '';
  guzik.dataset.uzbrojone = 'tak';
  guzik.textContent = 'Naciśnij ponownie';
  globalThis.setTimeout(() => {
    delete guzik.dataset.uzbrojone;
    guzik.textContent = napis;
  }, 5000);
  return false;
}

export function wypelnijPanel(korzen: Element, panelKod: string, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${panelKod} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

/*
zwiazPas wiąże nasłuch pasa z wykonawcą czynności.

Nasłuch stoi na korzeniu karty, a nie na samym pasie, bo pas bywa stawiany
przed dojściem odpowiedzi rdzenia i wtedy jeszcze go w drzewie nie ma.
*/
export function zwiazPas(
  korzen: Element,
  znacznik: string,
  wykonaj: (czynnosc: string) => void,
  przy: AddEventListenerOptions,
): void {
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const guzik = cel.closest<HTMLElement>(`[data-${znacznik}]`);
    const czynnosc = guzik?.getAttribute(`data-${znacznik}`) ?? '';
    if (czynnosc === '') return;
    zdarzenie.stopPropagation();
    wykonaj(czynnosc);
  }, przy);
}
