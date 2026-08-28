import { PROFILE_INZYNIERII } from './profile-inzynieria';
import { PROFILE_TRESCI } from './profile-tresc';
import { PROFILE_WIEDZY } from './profile-wiedza';
import { NARZEDZIA_WSPOLNE, POSTAC_DOMYSLNA, type ProfilModulu } from './profil-modulu';

/**
 * Rejestr profilów piętnastu modułów jest jedynym miejscem, w którym okno komunikacji szuka swojej postaci dla danego modułu, a moduł spoza rejestru dostaje profil wspólny zamiast blokady.
 */
const REJESTR: ReadonlyMap<string, ProfilModulu> = new Map(
  [...PROFILE_TRESCI, ...PROFILE_WIEDZY, ...PROFILE_INZYNIERII].map((p) => [p.kod, p]),
);

/** Profil modułu przypisany oknu komunikacji; moduł nierozpoznany przez rejestr dostaje profil wspólny zamiast blokady widoku. */
export function profilModulu(kod: string): ProfilModulu {
  return REJESTR.get(kod) ?? profilWspolny(kod);
}

/**
 * Profil okna pracującego w module spoza rejestru, z kodem pustym oznaczającym moduł jeszcze nieznany i wartościami domyślnymi postaci.
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
