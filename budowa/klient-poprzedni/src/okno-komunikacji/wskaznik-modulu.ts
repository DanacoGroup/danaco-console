import type { ProfilModulu } from './profil-modulu';

/** Wskaźnik modułu, w którym pracuje okno komunikacji. */
export interface WskaznikModulu {
  /** Element montowany w nagłówku okna. */
  element: HTMLElement;
  /** Przestawia wskaźnik na wskazany profil. */
  ustaw(profil: ProfilModulu): void;
}

/**
 * Wskaźnik bieżącego modułu.
 *
 * Operator ma w każdej chwili widzieć, w którym z piętnastu modułów pracuje
 * okno — bo od modułu zależy pasek promptu, panel akcji i kontekst. Wskaźnik
 * niesie słowo, nie samą barwę, a zmiana idzie przez `aria-live`, żeby czytnik
 * ekranu ogłosił przestawienie okna.
 */
export function utworzWskaznikModulu(): WskaznikModulu {
  const element = document.createElement('div');
  element.className = 'dc-modul';
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');

  const etykieta = document.createElement('span');
  etykieta.className = 'dc-modul__etykieta';
  etykieta.textContent = 'Moduł';

  const nazwa = document.createElement('span');
  nazwa.className = 'dc-modul__nazwa';

  const kod = document.createElement('code');
  kod.className = 'dc-modul__kod';

  element.append(etykieta, nazwa, kod);

  return {
    element,
    ustaw(profil) {
      element.dataset['modul'] = profil.kod;
      nazwa.textContent = profil.nazwa;
      kod.textContent = profil.kod;
      element.title = profil.przeznaczenie;
    },
  };
}
