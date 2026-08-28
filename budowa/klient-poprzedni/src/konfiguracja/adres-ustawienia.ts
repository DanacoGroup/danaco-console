import { ConfigAxis, ConfigScope, type ConfigEntry } from '../../../shared/contract';
import { nazwaOsi, nazwaZasiegu } from './zasiegi';

/**
 * Adres ustawienia — jedyny sposób wskazania miejsca w przestrzeni
 * konfiguracji. Przestrzeń ma dwa prostopadłe wymiary, więc adres ma cztery
 * człony: poziom zasięgu wraz z bytem poziomu oraz oś wraz z bytem osi.
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

/**
 * Adres otwierający okno: cała platforma na poziomie globalnym, czyli oba byty
 * puste. Od tego miejsca liczy się dziedziczenie, dopóki punkt widzenia nie
 * zostanie przestawiony.
 */
export function adresPoczatkowy(): AdresUstawienia {
  return {
    zasieg: ConfigScope.Global,
    bytZasiegu: '',
    os: ConfigAxis.Platform,
    bytOsi: '',
  };
}

/**
 * Adres, pod którym leży wpis konfiguracji; pola `scopeId` i `axisId` nieobecne
 * w odpowiedzi rdzenia dają byt pusty, a nieobecna oś znaczy oś platformy.
 */
export function adresWpisu(wpis: ConfigEntry): AdresUstawienia {
  return {
    zasieg: wpis.scope,
    bytZasiegu: wpis.scopeId ?? '',
    os: wpis.axis ?? ConfigAxis.Platform,
    bytOsi: wpis.axisId ?? '',
  };
}

/**
 * Czy dwa adresy wskazują to samo miejsce przestrzeni konfiguracji; przy pustym
 * bycie osi sama nazwa osi nie różnicuje, ponieważ poziom bez bytu jest jeden.
 */
export function tenSamAdres(pierwszy: AdresUstawienia, drugi: AdresUstawienia): boolean {
  return (
    pierwszy.zasieg === drugi.zasieg &&
    pierwszy.bytZasiegu === drugi.bytZasiegu &&
    pierwszy.bytOsi === drugi.bytOsi &&
    (pierwszy.bytOsi === '' || pierwszy.os === drugi.os)
  );
}

/**
 * Adres w jednym zdaniu: poziom zasięgu, byt poziomu, oś oraz byt osi, złączone
 * znakiem środkowej kropki. Człony puste do opisu nie wchodzą.
 */
export function opisAdresu(adres: AdresUstawienia): string {
  const czlony = [nazwaZasiegu(adres.zasieg)];
  if (adres.bytZasiegu !== '') czlony.push(adres.bytZasiegu);
  czlony.push(
    adres.bytOsi === '' ? nazwaOsi(adres.os) : `${nazwaOsi(adres.os)} ${adres.bytOsi}`,
  );
  return czlony.join(' · ');
}
