/** Pole wpisywania okna komunikacji, przyjmujące treść od operatora i wstawiane polecenia narzędzi modułu. */
export interface PoleWpisywania {
  element: HTMLElement;
  /** Ustawia ognisko na obszarze wpisywania. */
  ustawOgnisko(): void;
  /** Dokłada gotowe polecenie do treści; wysyłkę zatwierdza operator. */
  wstaw(tekst: string): void;
}

/**
 * Pole wpisywania z wysyłką przez Enter (Shift+Enter — nowy wiersz).
 *
 * Przycisk wysyłki pozostaje aktywny niezależnie od stanu połączenia;
 * przy rozłączeniu treść trafia do kolejki wychodzącej transportu.
 */
export function utworzPoleWpisywania(wyslij: (tresc: string) => void): PoleWpisywania {
  const element = document.createElement('form');
  element.className = 'dn-prompt';

  // Grot jest sygnaturą wejścia w pierwszej kolumnie siatki; czytany na głos nic nie niesie.
  const grot = document.createElement('span');
  grot.className = 'dn-prompt-grot';
  grot.textContent = '❯';
  grot.setAttribute('aria-hidden', 'true');

  const obszar = document.createElement('textarea');
  obszar.className = 'dn-prompt-obszar';
  obszar.rows = 3;
  obszar.placeholder = 'Napisz do modelu — Enter wysyła, Shift+Enter dodaje wiersz';
  obszar.setAttribute('aria-label', 'Treść wiadomości');

  // Wysyłka to jedyne działanie systemowe okna i jedyny wariant sygnałowy widoku.
  const przycisk = document.createElement('button');
  przycisk.type = 'submit';
  przycisk.className = 'dn-btn dn-btn--sygnal';
  przycisk.textContent = 'Wyślij';

  element.append(grot, obszar, przycisk);

  function przekaz(): void {
    const tresc = obszar.value.trim();
    if (tresc.length === 0) return;
    obszar.value = '';
    wyslij(tresc);
  }

  element.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    przekaz();
  });

  obszar.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter' && !zdarzenie.shiftKey) {
      zdarzenie.preventDefault();
      przekaz();
    }
  });

  return {
    element,
    ustawOgnisko: () => obszar.focus(),

    // Narzędzie modułu wstawia polecenie, nie wysyła go — turę otwiera operator, a treść zastana zostaje.
    wstaw(tekst) {
      const zastana = obszar.value;
      const rozdzielnik = zastana.length > 0 && !zastana.endsWith('\n') ? '\n' : '';
      obszar.value = `${zastana}${rozdzielnik}${tekst}`;
      obszar.focus();
      obszar.setSelectionRange(obszar.value.length, obszar.value.length);
    },
  };
}
