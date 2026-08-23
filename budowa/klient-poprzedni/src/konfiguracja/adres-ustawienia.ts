import { ConfigAxis, ConfigScope, type ConfigEntry } from '../../../shared/contract';
import { nazwaOsi, nazwaZasiegu } from './zasiegi';

/**
 * Adres ustawienia — jedyny sposób wskazania miejsca w przestrzeni
 * konfiguracji.
 *
 * Przestrzeń ma dwa prostopadłe wymiary, więc adres ma cztery człony: poziom
 * zasięgu wraz z bytem poziomu oraz oś wraz z bytem osi. Ten sam kształt służy
 * dwóm rolom, dlatego mieszka w jednym module:
 *
 *   • punkt widzenia — miejsce, z którego oglądana jest konfiguracja
 *     i względem którego liczone jest dziedziczenie;
 *   • adres zapisu — miejsce, w którym komenda `config.set` zapisze wartość.
 *
 * Byt pusty znaczy „poziom bez bytu": globalny na osi poziomów, platforma na
 * osi rozstrzygania.
 */
export interface AdresUstawienia {
  /** Poziom zasięgu. */
  zasieg: ConfigScope;
  /** Byt poziomu; pusty dla poziomu globalnego. */
  bytZasiegu: string;
  /** Oś rozstrzygania. */
  os: ConfigAxis;
  /** Byt osi; pusty dla osi platformy. */
  bytOsi: string;
}

/** Adres otwierający okno: cała platforma na poziomie globalnym. */
export function adresPoczatkowy(): AdresUstawienia {
  return {
    zasieg: ConfigScope.Global,
    bytZasiegu: '',
    os: ConfigAxis.Platform,
    bytOsi: '',
  };
}

/** Adres, pod którym leży wpis konfiguracji. */
export function adresWpisu(wpis: ConfigEntry): AdresUstawienia {
  return {
    zasieg: wpis.scope,
    bytZasiegu: wpis.scopeId ?? '',
    os: wpis.axis ?? ConfigAxis.Platform,
    bytOsi: wpis.axisId ?? '',
  };
}

/** Czy dwa adresy wskazują to samo miejsce przestrzeni konfiguracji. */
export function tenSamAdres(pierwszy: AdresUstawienia, drugi: AdresUstawienia): boolean {
  return (
    pierwszy.zasieg === drugi.zasieg &&
    pierwszy.bytZasiegu === drugi.bytZasiegu &&
    pierwszy.bytOsi === drugi.bytOsi &&
    (pierwszy.bytOsi === '' || pierwszy.os === drugi.os)
  );
}

/** Adres w jednym zdaniu: poziom, byt, oś, byt osi. */
export function opisAdresu(adres: AdresUstawienia): string {
  const czlony = [nazwaZasiegu(adres.zasieg)];
  if (adres.bytZasiegu !== '') czlony.push(adres.bytZasiegu);
  czlony.push(
    adres.bytOsi === '' ? nazwaOsi(adres.os) : `${nazwaOsi(adres.os)} ${adres.bytOsi}`,
  );
  return czlony.join(' · ');
}
