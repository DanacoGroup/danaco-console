import { liczbaOkienRozmowy, type ProfilModulu } from '../../okno-komunikacji/profil-modulu';
import { profilModulu } from '../../okno-komunikacji/rejestr-profilow';
import { KOD_MODULU } from './zrodlo-akcji-studio';

/**
 * Powiązanie modułu Studio z oknem rozmowy — opis, a nie drugie okno czatu.
 *
 * Chat Window jest jednym oknem na cały produkt i przy zmianie modułu przestawia
 * wygląd, narzędzia oraz kontekst (`okno-komunikacji/profil-modulu.ts`). Stoi na
 * scenie sesji, obok obszaru modułu, a nie wewnątrz niego, więc moduł z własnym
 * czatem prowadziłby drugi zapis tej samej rozmowy.
 *
 * Pasek nazywa to powiązanie i wylicza narzędzia promptu, które niesie profil
 * Studia: bez niego nie widać, że operacje Tools Panelu zleca się z okna
 * komunikacji sesji.
 *
 * Liczby okien rozmowy pasek nie podaje: profil Studia stoi na
 * `GRANICA_NIEPODANA`, więc wypisanie stamtąd liczby podawałoby wartość zastępczą
 * jako ustaloną.
 */
export interface PowiazanieRozmowy {
  element: HTMLElement;
}

/** Zdanie o postaci rozmowy modułu; każda postać ma własne, żadna nie milczy. */
function zdanieOPostaci(profil: ProfilModulu): string {
  if (liczbaOkienRozmowy(profil) === 0) {
    return (
      `Moduł ${profil.nazwa} nie prowadzi rozmowy — profil podaje postać „brak". ` +
      'Okna operacyjne poniżej są całą jego pracą.'
    );
  }
  if (profil.postacRozmowy === 'dymek-glosowy') {
    return (
      `Moduł ${profil.nazwa} rozmawia dymkiem głosowym, a nie oknem czatu. ` +
      'Dopóki rdzeń nie niesie kanału głosowego, dymek ma to powiedzieć wprost.'
    );
  }
  return (
    `Rozmowa modułu ${profil.nazwa} toczy się w oknie komunikacji sesji — tym samym, które pas ` +
    'osadzenia wskazuje wyżej jako windowId. Operacje zleca się stamtąd, wierszem polecenia przy ' +
    'kursorze albo z Tools Panelu, a wynik wchodzi WPROST DO TREŚCI dokumentu jako zmiana ' +
    'oznaczona autorem „model" — kanwy tekstowej, na której stawał osobno, już nie ma.'
  );
}

/** Zdanie o pamięci wątku — moduł gubiący wątek po cichu wygląda jak awaria. */
function zdanieOPamieci(profil: ProfilModulu): string {
  return profil.pamiecSesyjna
    ? 'Wątek rozmowy przeżywa zmianę modułu i zamknięcie okna operacyjnego: profil ma pamięć sesyjną.'
    : 'Wątek rozmowy ginie przy zamknięciu okna — profil nie ma pamięci sesyjnej. Zapisz, co ma zostać.';
}

export function utworzPowiazanieRozmowy(): PowiazanieRozmowy {
  const profil = profilModulu(KOD_MODULU);

  const element = document.createElement('section');
  element.className = 'ms-rozmowa';
  element.dataset['postac'] = profil.postacRozmowy;
  element.setAttribute('aria-label', 'Powiązanie modułu Studio z oknem rozmowy');

  const tytul = document.createElement('p');
  tytul.className = 'ms-rozmowa__tytul';
  tytul.textContent = 'Okno komunikacji modułu';

  const postac = document.createElement('p');
  postac.className = 'dn-pole-opis';
  postac.textContent = zdanieOPostaci(profil);

  const pamiec = document.createElement('p');
  pamiec.className = 'dn-pole-opis';
  pamiec.textContent = zdanieOPamieci(profil);

  element.append(tytul, postac, pamiec);

  // Moduł bez rozmowy nie dostaje wykazu narzędzi promptu: byłby to wykaz
  // czynności, których nie ma gdzie wykonać.
  if (liczbaOkienRozmowy(profil) === 0) return { element };

  const wstep = document.createElement('p');
  wstep.className = 'dn-pole-opis';
  wstep.textContent =
    'Pasek narzędzi promptu tego okna niesie w module Studio pozycje wymienione niżej. Każda ' +
    'wstawia gotowe polecenie do pola wypowiedzi, a wykonuje je model wywołaniem narzędzia ' +
    'platformy — nie jest to osobna komenda kontraktu.';

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-rozmowa__wykaz';
  for (const narzedzie of profil.narzedzia) {
    const pozycja = document.createElement('li');
    const nazwa = document.createElement('strong');
    nazwa.textContent = narzedzie.nazwa;
    const opis = document.createElement('span');
    opis.textContent = ` — ${narzedzie.opis}`;
    pozycja.dataset['narzedzie'] = narzedzie.kod;
    pozycja.append(nazwa, opis);
    wykaz.append(pozycja);
  }

  element.append(wstep, wykaz);
  return { element };
}
