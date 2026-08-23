import { PROFILE_INZYNIERII } from './profile-inzynieria';
import { PROFILE_TRESCI } from './profile-tresc';
import { PROFILE_WIEDZY } from './profile-wiedza';
import { NARZEDZIA_WSPOLNE, POSTAC_DOMYSLNA, type ProfilModulu } from './profil-modulu';

/**
 * Rejestr profilów piętnastu modułów — jedyne miejsce, w którym okno
 * komunikacji szuka swojej postaci dla danego modułu.
 *
 * Wykaz jest informacyjny, nie bramą: moduł spoza rejestru nie
 * blokuje okna ani nie wywraca widoku. Dostaje profil wspólny — wskaźnik
 * niesie jego kod, pasek niesie arsenał wspólny Chat Window, a panel akcji
 * i tak bierze zawartość z katalogu rdzenia po kodzie modułu.
 */
const REJESTR: ReadonlyMap<string, ProfilModulu> = new Map(
  [...PROFILE_TRESCI, ...PROFILE_WIEDZY, ...PROFILE_INZYNIERII].map((p) => [p.kod, p]),
);

/** Profil modułu; moduł nierozpoznany dostaje profil wspólny. */
export function profilModulu(kod: string): ProfilModulu {
  return REJESTR.get(kod) ?? profilWspolny(kod);
}

/**
 * Profil okna pracującego w module spoza rejestru.
 *
 * Kod pusty znaczy „moduł jeszcze nieznany" — okno czeka na odpowiedź rdzenia
 * i mówi to wprost, zamiast udawać moduł, którego nie zna.
 *
 * Postać rozmowy, granica okien i pamięć sesyjna biorą wartości `POSTAC_DOMYSLNA`,
 * bo dla modułu spoza rejestru nie ma ustalonej postaci.
 */
function profilWspolny(kod: string): ProfilModulu {
  const nazwa = kod.length > 0 ? kod : 'moduł nieustalony';
  const przeznaczenie =
    kod.length > 0
      ? 'Moduł spoza rejestru profilów klienta — okno pracuje na arsenale wspólnym'
      : 'Moduł okna nie został jeszcze potwierdzony przez rdzeń';
  return {
    kod,
    nazwa,
    przeznaczenie,
    oknaKontekstu: [],
    narzedzia: NARZEDZIA_WSPOLNE,
    postacRozmowy: POSTAC_DOMYSLNA.postacRozmowy,
    granicaOkien: POSTAC_DOMYSLNA.granicaOkien,
    pamiecSesyjna: POSTAC_DOMYSLNA.pamiecSesyjna,
    oknaObowiazkowe: [],
  };
}
