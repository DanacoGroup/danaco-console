/* Składnik biblioteki rysuje bryłę modułu: sześcian programu 3×3×3 składany z klocków warstwami od dołu, z błyskiem krawędzi przy każdym zatrzasku, pierścieniem po ukończonej warstwie i poświatą rdzenia po złożeniu.
   BRYŁA MODUŁU — składnik biblioteki

   Sześcian programu 3×3×3 składany z klocków warstwami od dołu. Każdy zatrzask
   daje błysk krawędzi, ukończona warstwa puszcza pierścień po podłożu, a wnętrze
   rozświetla się rdzeniem. Po złożeniu moduł pracuje, po chwili rozkłada się
   i cykl rusza od nowa.

   Okno wstawia puste pole i nic więcej:

       <div class="dn-bryla" data-bryla
            role="img" aria-label="Budowa modułu Danaco Console"></div>

   Bryła NIE MA własnego tła — kładzie się wprost na powierzchni okna. Jedna
   paleta czyta się i na atramencie, i na papierze, więc scena nie zmienia się
   z motywem. Kadr wpisuje się w pole przy każdej zmianie rozmiaru, więc nic nie
   jest ucinane i nic nie wymaga wygaszania krawędzi.

   Barwy wyłącznie z żetonów `--dn-bryla-*`. Okno nie podaje ani jednej wartości
   wizualnej — podaje wyłącznie stan.

   Sterowanie z okna:
       var b = document.querySelector('[data-bryla]').bryla;
       b.postep(0.42);        // 0..1 — klocki układają się wraz z postępem
       b.pokaz();             // powrót do pętli automatycznej
       b.pracuje(true);       // praca w toku — żywsze tempo
       b.domknij();           // błysk domknięcia i pierścień
       b.etap('Kopiowanie plików');            // podpis nad modułem
       b.etykiety(['RDZEŃ', 'MODUŁY', ...]);   // nazwy warstw
       b.tempo(1.2);
       b.wstrzymaj(); b.wznow(); b.zdejmij();

   Atrybuty pola:
       data-warstwy="RDZEŃ, MODUŁY, ŚRODOWISKO, INTERFEJS"
       data-postep="0.4"        postęp startowy — wyłącza pętlę
       data-montaz="9"          czas montażu w sekundach

   Podpis etapu jest domyślnie wyłączony i mieści się dopiero w polu szerszym
   niż 420 px. Przy `prefers-reduced-motion` ruch kamery jest ograniczony.
   ========================================================================== */
(function () {
'use strict';

var ETYKIETY = ['RDZEŃ', 'MODUŁY', 'ŚRODOWISKO', 'INTERFEJS'];
var N = 3, SZ = 0.52, ODSTEP = 0.025, KROK = SZ + ODSTEP;

function ogranicz(v, a, b) { return v < a ? a : (v > b ? b : v); }
function losowo(a, b) { return a + Math.random() * (b - a); }
function zwolnienie(t) { return 1 - Math.pow(1 - ogranicz(t, 0, 1), 3); }
function wygladzenie(t) { t = ogranicz(t, 0, 1); return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2; }
function miedzy(a, b, t) { return a + (b - a) * t; }

/* ————————————————————————————————————————————————————————— barwy z żetonów
   Żeton bywa zapisany szesnastkowo, funkcją `rgb()`, `rgba()` albo odwołaniem
   do innego żetonu. Zamiast rozbierać każdy zapis z osobna, wartość idzie do
   próbki i wraca znormalizowana przez przeglądarkę — razem z kryciem, bo ściany
   klocków są półprzezroczyste i krycie należy do żetonu, nie do rysunku. */
/* Czytnik barw stoi w `narzedzia-okien.js` — potrzebuje go każdy składnik
   rysujący na płótnie. */
var czytnikBarw = window.DanacoNarzedzia.czytnikBarw;
function kry(c, a) { return 'rgba(' + c.rgb[0] + ',' + c.rgb[1] + ',' + c.rgb[2] + ',' + (a * c.a) + ')'; }

function zbudujPalete(host) {
  var cz = czytnikBarw(host);
  var P = {
    linia:       cz.barwa('--dn-bryla-linia', '#68A8F6'),
    poswiata:    cz.barwa('--dn-bryla-poswiata', '#3474C8'),
    lagodna:     cz.barwa('--dn-bryla-lagodna', '#84B6F0'),
    wezel:       cz.barwa('--dn-bryla-wezel', '#A8CEF6'),
    blysk:       cz.barwa('--dn-bryla-blysk', '#F0F8FF'),
    blat:        cz.barwa('--dn-bryla-blat', 'rgba(41,66,101,.97)'),
    blatWierzch: cz.barwa('--dn-bryla-blat-wierzch', 'rgba(52,82,122,.97)'),
    bok:         cz.barwa('--dn-bryla-bok', 'rgba(26,45,71,.97)'),
    bokCien:     cz.barwa('--dn-bryla-bok-cien', 'rgba(17,30,48,.97)'),
    cien:        cz.barwa('--dn-bryla-cien', 'rgba(16,36,66,.26)'),
    plakietka:   cz.barwa('--dn-bryla-plakietka', 'rgba(22,36,56,.92)'),
    plakietkaTekst: cz.barwa('--dn-bryla-plakietka-tekst', '#E2EEFC')
  };
  cz.zdejmij();
  return P;
}

/* Instancja bryły: montaż warstw sześcianu, cykl pracy i rozkładu oraz reakcja na zmianę rozmiaru pola i widoczność w oknie przeglądarki. */
function zaloz(host, opcje) {
  if (!host || host.bryla) return host ? host.bryla : null;
  opcje = opcje || {};

  var plotno = document.createElement('canvas');
  plotno.setAttribute('aria-hidden', 'true');
  host.appendChild(plotno);
  var ctx = plotno.getContext('2d');

  var pytanieRuchu = window.matchMedia('(prefers-reduced-motion: reduce)');
  var SKROMNIE = pytanieRuchu.matches;
  var P = zbudujPalete(host);

  var W = 0, H = 0, DPR = 1, jednostka = 1, sX = 0, sY = 0;
  var kam = { obrot: 0.62, pochyl: 0.46, odleglosc: 6.0, ognisko: 3.0 };

  function rzut(x, y, z) {
    var co = Math.cos(kam.obrot), so = Math.sin(kam.obrot);
    var rx = x * co - z * so, rz = x * so + z * co;
    var cp = Math.cos(kam.pochyl), sp = Math.sin(kam.pochyl);
    var yc = y * cp + rz * sp, zc = -y * sp + rz * cp + kam.odleglosc;
    var s = kam.ognisko * jednostka / Math.max(0.08, zc);
    return { x: sX + rx * s, y: sY - yc * s, s: s, d: zc };
  }
  function obrys(pts) {
    ctx.beginPath(); ctx.moveTo(pts[0].x, pts[0].y);
    for (var i = 1; i < pts.length; i++) ctx.lineTo(pts[i].x, pts[i].y);
    ctx.closePath();
  }

  var klocki = [], pierscienie = [], iskry = [], etykiety = ETYKIETY.slice();

  function obwod() {
    var t = [], i;
    for (i = 0; i < 3; i++) {
      var x = losowo(-0.32, 0.32), z = losowo(-0.32, 0.32);
      var pts = [{ x: x, z: z }], os = Math.random() < 0.5, k;
      for (k = 0; k < 2; k++) {
        var d = losowo(0.12, 0.28) * (Math.random() < 0.5 ? -1 : 1);
        if (os) x = ogranicz(x + d, -0.34, 0.34); else z = ogranicz(z + d, -0.34, 0.34);
        pts.push({ x: x, z: z }); os = !os;
      }
      t.push(pts);
    }
    return t;
  }
  /* Kolejność montażu: warstwami od dołu, wewnątrz warstwy od środka na zewnątrz
     — moduł rośnie tak, jak się go buduje, a nie w przypadkowej kolejności. */
  function zbudujKlocki() {
    klocki = []; pierscienie = []; iskry = [];
    var x, y, z, kolejnosc = [];
    for (y = 0; y < N; y++) for (x = 0; x < N; x++) for (z = 0; z < N; z++) kolejnosc.push([x, y, z]);
    kolejnosc.sort(function (a, b) {
      if (a[1] !== b[1]) return a[1] - b[1];
      var ra = Math.hypot(a[0] - 1, a[2] - 1), rb = Math.hypot(b[0] - 1, b[2] - 1);
      if (ra !== rb) return ra - rb;
      return Math.atan2(a[2] - 1, a[0] - 1) - Math.atan2(b[2] - 1, b[0] - 1);
    });
    for (var i = 0; i < kolejnosc.length; i++) {
      var o = kolejnosc[i], a = losowo(0, Math.PI * 2), r = losowo(1.85, 2.55);
      klocki.push({
        gx: o[0], gy: o[1], gz: o[2],
        x: (o[0] - 1) * KROK, y: (o[1] - 1) * KROK, z: (o[2] - 1) * KROK,
        skad: { x: Math.cos(a) * r, y: losowo(1.15, 2.05), z: Math.sin(a) * r },
        skret: losowo(-1.4, 1.4), kolej: i / kolejnosc.length,
        widok: 0, blysk: 0, unos: 0, wierzch: (o[1] === N - 1),
        obwod: (o[1] === N - 1) ? obwod() : null
      });
    }
  }
  zbudujKlocki();

  var SCIANY = [[0, 1, 2, 3], [4, 5, 6, 7], [0, 1, 5, 4], [1, 2, 6, 5], [2, 3, 7, 6], [3, 0, 4, 7]];
  function wierzcholki(cx, cy, cz, s, skret) {
    var h = s / 2, v = [], zn = [[-1, 1, -1], [1, 1, -1], [1, 1, 1], [-1, 1, 1], [-1, -1, -1], [1, -1, -1], [1, -1, 1], [-1, -1, 1]];
    var c = Math.cos(skret || 0), sn = Math.sin(skret || 0);
    for (var i = 0; i < 8; i++) {
      var x = zn[i][0] * h, y = zn[i][1] * h, z = zn[i][2] * h;
      v.push(rzut(cx + (x * c - z * sn), cy + y, cz + (x * sn + z * c)));
    }
    return v;
  }
  function rysujKlocek(C, alfa, t) {
    var k = zwolnienie(C.widok);
    var x = miedzy(C.skad.x, C.x, k), y = miedzy(C.skad.y, C.y, k) + C.unos * 2.6, z = miedzy(C.skad.z, C.z, k);
    var v = wierzcholki(x, y, z, SZ, C.skret * (1 - k)), i, sciany = [];
    for (i = 0; i < SCIANY.length; i++) {
      var f = SCIANY[i];
      sciany.push({ f: f, i: i, d: (v[f[0]].d + v[f[1]].d + v[f[2]].d + v[f[3]].d) / 4 });
    }
    sciany.sort(function (a, b) { return b.d - a.d; });
    ctx.save(); ctx.globalAlpha = alfa;
    for (i = 0; i < sciany.length; i++) {
      var S = sciany[i];
      obrys([v[S.f[0]], v[S.f[1]], v[S.f[2]], v[S.f[3]]]);
      ctx.fillStyle = S.i === 0 ? (C.wierzch ? P.blatWierzch.css : P.blat.css)
                    : (S.i === 1 ? P.bokCien.css : (S.i % 2 ? P.bok.css : P.bokCien.css));
      ctx.fill();
      ctx.strokeStyle = kry(P.linia, (S.i === 0 ? 0.9 : 0.5) + 0.4 * C.blysk);
      ctx.lineWidth = 1 + 1.7 * C.blysk;
      ctx.stroke();
    }
    /* Obwód na blacie najwyższej warstwy — znak, że moduł pracuje. */
    if (C.wierzch && C.obwod && k > 0.9) {
      ctx.save();
      obrys([v[0], v[1], v[2], v[3]]); ctx.clip();
      var ca = (k - 0.9) / 0.1;
      for (i = 0; i < C.obwod.length; i++) {
        var tr = C.obwod[i], j;
        ctx.strokeStyle = kry(P.lagodna, 0.55 * ca);
        ctx.lineWidth = Math.max(0.7, 0.012 * v[0].s);
        ctx.beginPath();
        for (j = 0; j < tr.length; j++) {
          var pp = rzut(x + tr[j].x * SZ, y + SZ / 2 + 0.001, z + tr[j].z * SZ);
          if (j === 0) ctx.moveTo(pp.x, pp.y); else ctx.lineTo(pp.x, pp.y);
        }
        ctx.stroke();
        var f2 = (t * 0.35 + i * 0.31 + C.kolej) % 1;
        var seg = Math.min(tr.length - 2, Math.floor(f2 * (tr.length - 1)));
        var lf = f2 * (tr.length - 1) - seg;
        var A = tr[seg], B = tr[seg + 1];
        var np = rzut(x + miedzy(A.x, B.x, lf) * SZ, y + SZ / 2 + 0.002, z + miedzy(A.z, B.z, lf) * SZ);
        var ds = Math.max(1.5, 0.022 * np.s);
        ctx.fillStyle = kry(P.blysk, 0.9 * ca);
        ctx.shadowBlur = 8; ctx.shadowColor = kry(P.poswiata, 0.95);
        ctx.fillRect(np.x - ds / 2, np.y - ds / 2, ds, ds);
        ctx.shadowBlur = 0;
      }
      ctx.restore();
    }
    if (C.blysk > 0.02) {
      obrys([v[0], v[1], v[2], v[3]]);
      ctx.fillStyle = kry(P.blysk, 0.3 * C.blysk); ctx.fill();
    }
    ctx.restore();
  }

  /* ——————————————————————————————————————————————————————————————— kadr
     Bryła nie ma tła, więc nic nie zasłoni miejsca, w którym rysunek wyszedłby
     poza płótno. Kadr wpisuje się w pole pomiarem: skrajne punkty modułu liczone
     dla kilku położeń kamery muszą zmieścić się w ramce z zapasem. */
  function punktyKadru() {
    var a = [], i, R = (N * KROK) * 0.78;
    for (i = 0; i < 20; i++) {
      var t = i / 20 * Math.PI * 2;
      a.push([Math.cos(t) * R, -(N * KROK) / 2 - 0.28, Math.sin(t) * R]);
      a.push([Math.cos(t) * R, (N * KROK) / 2 + 0.34, Math.sin(t) * R]);
    }
    return a;
  }
  function wpiszKadr() {
    var obrot0 = kam.obrot, pochyl0 = kam.pochyl;
    var zapas = Math.max(6, Math.min(W, H) * 0.04);
    var zapasP = (podpis && W > 420) ? Math.max(zapas, Math.min(W * 0.26, 180)) : zapas;
    /* Górna część pola bywa zajęta przez treść leżącą NAD sceną — moduł składa
       się wtedy niżej, a zwolniona góra zostaje miejscem dla klocków w locie. */
    var x0 = zapas, x1 = Math.max(x0 + 40, W - zapasP);
    var y0 = zapas + H * kadrGora, y1 = Math.max(y0 + 40, H - zapas);
    var szer = x1 - x0, wys = y1 - y0, pts = punktyKadru(), i;
    kam.obrot = 0;
    for (var it = 0; it < 3; it++) {
      var pochylenia = [0.41, 0.46, 0.51], minx = 1e9, maxx = -1e9, miny = 1e9, maxy = -1e9;
      for (var w = 0; w < pochylenia.length; w++) {
        kam.pochyl = pochylenia[w];
        for (i = 0; i < pts.length; i++) {
          var q = rzut(pts[i][0], pts[i][1], pts[i][2]);
          if (q.x < minx) minx = q.x;
          if (q.x > maxx) maxx = q.x;
          if (q.y < miny) miny = q.y;
          if (q.y > maxy) maxy = q.y;
        }
      }
      var k = Math.min(szer / Math.max(1, maxx - minx), wys / Math.max(1, maxy - miny));
      k = Math.max(0.3, Math.min(1.7, k));
      jednostka *= k;
      var mx = (minx + maxx) / 2, my = (miny + maxy) / 2;
      sX += (x0 + szer / 2) - (sX + (mx - sX) * k);
      sY += (y0 + wys / 2) - (sY + (my - sY) * k);
    }
    kam.obrot = obrot0; kam.pochyl = pochyl0;
  }
  function wymiar() {
    DPR = Math.min(2, window.devicePixelRatio || 1);
    W = host.clientWidth || 260;
    H = host.clientHeight || 200;
    plotno.width = Math.max(1, Math.round(W * DPR));
    plotno.height = Math.max(1, Math.round(H * DPR));
    plotno.style.width = W + 'px';
    plotno.style.height = H + 'px';
    ctx.setTransform(DPR, 0, 0, DPR, 0, 0);
    jednostka = Math.min(W, H * 1.35) * 0.42;
    sX = W * 0.5; sY = H * 0.5;
    wpiszKadr();
  }

  /* ——————————————————————————————————————————————————————————————— cykl */
  var MONTAZ = Number(host.dataset.montaz || opcje.montaz || 9) || 9;
  var CYKL = { trwanie: 3.2, rozklad: 1.6, przerwa: 0.5 };
  var faza = 'montaz', fT = 0, postep = 0, rozklad = 0, zewnetrzny = false;
  var blyskKonca = 0, biegnie = true, ostatniTS = 0, tGlob = 0, skan = -1;
  var tempo = 1, podpis = false, tekstEtapu = '', klatkaId = 0;
  var kadrGora = ogranicz(parseFloat(host.dataset.kadrGora || opcje.kadrGora || 0) || 0, 0, 0.8);
  var PRZYGASZENIE = 0.62;

  if (host.dataset.warstwy) {
    etykiety = host.dataset.warstwy.split(',').map(function (s) { return s.trim(); });
  }
  if (host.hasAttribute('data-postep')) {
    zewnetrzny = true;
    postep = ogranicz(parseFloat(host.dataset.postep) || 0, 0, 1);
  }

  function odNowa() {
    faza = 'montaz'; fT = 0; postep = 0; rozklad = 0; blyskKonca = 0; skan = -1;
    for (var i = 0; i < klocki.length; i++) { klocki[i].widok = 0; klocki[i].blysk = 0; klocki[i].unos = 0; }
  }
  function warstwaZlozona(gy) {
    for (var i = 0; i < klocki.length; i++) if (klocki[i].gy === gy && klocki[i].widok < 0.9) return false;
    return true;
  }
  function przelicz(dt) {
    tGlob += dt;
    if (!zewnetrzny) {
      fT += dt;
      if (faza === 'montaz') {
        postep = ogranicz(fT / MONTAZ, 0, 1);
        if (fT >= MONTAZ) {
          faza = 'praca'; fT = 0; blyskKonca = 1;
          pierscienie.push({ y: -(N * KROK) / 2, t: 0, zycie: 1.6, w: 1.8 });
        }
      } else if (faza === 'praca') {
        postep = 1;
        if (fT >= CYKL.trwanie) { faza = 'rozklad'; fT = 0; }
      } else if (faza === 'rozklad') {
        rozklad = ogranicz(fT / CYKL.rozklad, 0, 1);
        for (var q = 0; q < klocki.length; q++) {
          var C2 = klocki[q];
          C2.unos = wygladzenie(ogranicz((rozklad - (1 - C2.kolej) * 0.45) / 0.55, 0, 1));
        }
        if (fT >= CYKL.rozklad) { faza = 'przerwa'; fT = 0; }
      } else if (faza === 'przerwa') {
        if (fT >= CYKL.przerwa) odNowa();
      }
    }
    kam.obrot = 0.62 + Math.sin(tGlob * (SKROMNIE ? 0.06 : 0.15)) * 0.30 + tGlob * (SKROMNIE ? 0 : 0.05);
    kam.pochyl = 0.46 + Math.sin(tGlob * 0.12) * 0.035;

    for (var i = 0; i < klocki.length; i++) {
      var C = klocki[i];
      var cel = ogranicz((postep - C.kolej * 0.93) / 0.09, 0, 1), przed = C.widok;
      C.widok += (cel - C.widok) * Math.min(1, dt * (cel > C.widok ? 5.5 : 7));
      if (przed < 0.94 && C.widok >= 0.94) {
        C.blysk = 1;
        iskry.push({ x: C.x, y: C.y, z: C.z, t: 0, zycie: 0.55 });
        if (warstwaZlozona(C.gy)) pierscienie.push({ y: C.y - KROK / 2, t: 0, zycie: 1.3, w: 1.1 });
      }
      C.blysk = Math.max(0, C.blysk - dt * 2.4);
    }
    for (i = pierscienie.length - 1; i >= 0; i--) {
      pierscienie[i].t += dt;
      if (pierscienie[i].t >= pierscienie[i].zycie) pierscienie.splice(i, 1);
    }
    for (i = iskry.length - 1; i >= 0; i--) {
      iskry[i].t += dt;
      if (iskry[i].t >= iskry[i].zycie) iskry.splice(i, 1);
    }
    blyskKonca = Math.max(0, blyskKonca - dt * 1.2);
    if (postep >= 0.999) {
      skan += dt * 0.85;
      if (skan > (N * KROK) / 2 + 0.3) skan = -(N * KROK) / 2 - 0.3;
    }
  }

  function rysujCien(al) {
    var c = rzut(0, -(N * KROK) / 2 - 0.02, 0);
    ctx.save();
    ctx.globalAlpha = al * ogranicz(postep * 1.5, 0, 1);
    var g = ctx.createRadialGradient(c.x, c.y, 1, c.x, c.y, jednostka * 0.62);
    g.addColorStop(0, P.cien.css);
    g.addColorStop(1, kry(P.cien, 0));
    ctx.fillStyle = g; ctx.fillRect(0, 0, W, H);
    ctx.restore();
  }
  function rysujPierscienie(al) {
    ctx.save();
    for (var r = 0; r < pierscienie.length; r++) {
      var R = pierscienie[r], k = R.t / R.zycie;
      ctx.globalAlpha = al * (1 - k) * (1 - k) * 0.65;
      ctx.strokeStyle = kry(P.wezel, 0.9);
      ctx.lineWidth = (2.2 * (1 - k) + 0.4) * R.w;
      var promien = 0.55 + k * 1.9;
      ctx.beginPath();
      for (var i = 0; i <= 44; i++) {
        var a = i / 44 * Math.PI * 2, p = rzut(Math.cos(a) * promien, R.y, Math.sin(a) * promien);
        if (i === 0) ctx.moveTo(p.x, p.y); else ctx.lineTo(p.x, p.y);
      }
      ctx.stroke();
    }
    ctx.restore();
  }
  function rysujRdzen(al) {
    var v = ogranicz((postep - 0.3) / 0.55, 0, 1) * (0.75 + 0.25 * Math.sin(tGlob * 2.2));
    if (v <= 0.01) return;
    var c = rzut(0, 0, 0);
    ctx.save(); ctx.globalAlpha = al;
    var g = ctx.createRadialGradient(c.x, c.y, 0, c.x, c.y, jednostka * 0.46);
    g.addColorStop(0, kry(P.lagodna, 0.3 * v + 0.22 * blyskKonca));
    g.addColorStop(1, kry(P.poswiata, 0));
    ctx.fillStyle = g; ctx.fillRect(0, 0, W, H);
    ctx.restore();
  }
  function rysujIskry(al) {
    ctx.save();
    for (var i = 0; i < iskry.length; i++) {
      var S = iskry[i], k = S.t / S.zycie, p = rzut(S.x, S.y, S.z);
      ctx.globalAlpha = al * (1 - k) * (1 - k) * 0.55;
      ctx.strokeStyle = kry(P.blysk, 0.9);
      ctx.lineWidth = 1.4 * (1 - k) + 0.4;
      ctx.beginPath();
      ctx.arc(p.x, p.y, (0.3 + k * 0.26) * p.s, 0, Math.PI * 2);
      ctx.stroke();
    }
    ctx.restore();
  }
  /* Podpis etapu mieści się dopiero w polu szerszym niż 420 px — w węższym
     odnośnik i plakietka nachodziłyby na moduł. Wtedy nazwę niesie `aria-label`. */
  function rysujPodpis(al) {
    if (!podpis || W < 420 || al <= 0.05) return;
    var txt = tekstEtapu || etykiety[Math.min(etykiety.length - 1, Math.floor(postep * etykiety.length))];
    if (!txt) return;
    var an = rzut(0.95, (N * KROK) / 2 - 0.1, 0.95);
    var styl = getComputedStyle(host);
    ctx.save();
    ctx.font = '11px ' + (styl.getPropertyValue('--dn-ff-mono').trim() || 'ui-monospace, monospace');
    if ('letterSpacing' in ctx) ctx.letterSpacing = '1.4px';
    ctx.textBaseline = 'middle';
    var tw = ctx.measureText(txt).width;
    var kolX = Math.max(Math.min(W - 26 - tw, an.x + Math.min(70, W * 0.08)), an.x + 22);
    var ty = ogranicz(an.y, 18, H - 18);
    ctx.globalAlpha = al;
    ctx.strokeStyle = kry(P.linia, 0.75); ctx.lineWidth = 1.1;
    ctx.beginPath();
    ctx.moveTo(an.x + 7, an.y); ctx.lineTo(an.x + 16, ty); ctx.lineTo(kolX - 12, ty);
    ctx.stroke();
    ctx.fillStyle = kry(P.wezel, 1); ctx.fillRect(kolX - 15, ty - 1.5, 3, 3);
    ctx.fillStyle = P.plakietka.css; ctx.strokeStyle = kry(P.linia, 0.55); ctx.lineWidth = 1;
    ctx.beginPath();
    if (ctx.roundRect) ctx.roundRect(kolX - 7, ty - 9, tw + 14, 18, 4);
    else ctx.rect(kolX - 7, ty - 9, tw + 14, 18);
    ctx.fill(); ctx.stroke();
    ctx.fillStyle = kry(P.plakietkaTekst, 1);
    ctx.fillText(txt, kolX, ty);
    ctx.restore();
  }
  function rysuj() {
    ctx.clearRect(0, 0, W, H);
    var al = zewnetrzny ? 1
      : (faza === 'przerwa' ? 0
      : (faza === 'rozklad' ? 1 - wygladzenie(ogranicz((rozklad - 0.2) / 0.8, 0, 1)) : 1));
    if (al <= 0.01) return;
    rysujCien(al);
    rysujPierscienie(al);
    rysujRdzen(al);
    var lista = [], i;
    for (i = 0; i < klocki.length; i++) if (klocki[i].widok > 0.01) lista.push(klocki[i]);
    lista.sort(function (a, b) {
      var ka = zwolnienie(a.widok), kb = zwolnienie(b.widok);
      var da = rzut(miedzy(a.skad.x, a.x, ka), miedzy(a.skad.y, a.y, ka) + a.unos * 2.6, miedzy(a.skad.z, a.z, ka)).d;
      var db = rzut(miedzy(b.skad.x, b.x, kb), miedzy(b.skad.y, b.y, kb) + b.unos * 2.6, miedzy(b.skad.z, b.z, kb)).d;
      return db - da;
    });
    for (i = 0; i < lista.length; i++) {
      rysujKlocek(lista[i], al * ogranicz(lista[i].widok * 1.4, 0, 1) * (1 - lista[i].unos * 0.9), tGlob);
    }
    rysujIskry(al);
    rysujPodpis(al * 0.95);
    przygasPodTekstem();
  }

  /* Górny pas pola leży pod treścią okna — scena jest tam tłem, nie tematem,
     więc gaśnie do połowy siły. Nie jest to winieta: nic nie dorysowujemy,
     tylko zdejmujemy krycie samemu rysunkowi, i wyłącznie tam, gdzie nachodzi
     na tekst. Klocki w locie zostają widoczne, ale przestają z nim walczyć. */
  function przygasPodTekstem() {
    if (kadrGora <= 0) return;
    var dol = H * kadrGora;
    var g = ctx.createLinearGradient(0, 0, 0, dol);
    g.addColorStop(0, 'rgba(0,0,0,' + PRZYGASZENIE + ')');
    g.addColorStop(1, 'rgba(0,0,0,0)');
    ctx.save();
    ctx.globalCompositeOperation = 'destination-out';
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, W, dol);
    ctx.restore();
  }
  function klatka(ts) {
    if (!biegnie) return;
    if (!ostatniTS) ostatniTS = ts;
    var dt = Math.min(0.05, (ts - ostatniTS) / 1000) * tempo;
    ostatniTS = ts;
    przelicz(dt); rysuj();
    klatkaId = requestAnimationFrame(klatka);
  }

  /* ——————————————————————————————————————————— reakcja na otoczenie */
  function przemaluj() { P = zbudujPalete(host); }
  var obserwatorMotywu = new MutationObserver(przemaluj);
  obserwatorMotywu.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });
  var pytanieMotywu = window.matchMedia('(prefers-color-scheme: dark)');
  function naZmianeRuchu() { SKROMNIE = pytanieRuchu.matches; }
  if (pytanieMotywu.addEventListener) {
    pytanieMotywu.addEventListener('change', przemaluj);
    pytanieRuchu.addEventListener('change', naZmianeRuchu);
  }
  var obserwatorRozmiaru = window.ResizeObserver ? new ResizeObserver(wymiar) : null;
  if (obserwatorRozmiaru) obserwatorRozmiaru.observe(host);
  else window.addEventListener('resize', wymiar);

  /* Bryła poza widokiem nie ma komu nic pokazywać — pętla staje, zamiast palić
     klatki pod ukrytym krokiem kreatora. */
  var obserwatorWidoku = window.IntersectionObserver
    ? new IntersectionObserver(function (wpisy) {
        if (wpisy[wpisy.length - 1].isIntersecting) api.wznow(); else api.wstrzymaj();
      })
    : null;
  if (obserwatorWidoku) obserwatorWidoku.observe(host);
  function naWidocznosc() { if (document.hidden) api.wstrzymaj(); else api.wznow(); }
  document.addEventListener('visibilitychange', naWidocznosc);

  var api = {
    postep: function (p) {
      zewnetrzny = true;
      postep = ogranicz(Number(p) || 0, 0, 1);
      for (var i = 0; i < klocki.length; i++) klocki[i].unos = 0;
      return postep;
    },
    /* Powrót do pętli automatycznej — montaż, praca, rozkład, od nowa. */
    pokaz: function () { zewnetrzny = false; odNowa(); },
    /* Praca w toku: moduł składa się żwawiej. Nie zmienia treści sceny — sam rytm. */
    pracuje: function (wl) { tempo = wl ? 1.35 : 1; },
    domknij: function () {
      blyskKonca = 1;
      pierscienie.push({ y: -(N * KROK) / 2, t: 0, zycie: 1.6, w: 1.8 });
    },
    etap: function (t) { tekstEtapu = t || ''; podpis = !!tekstEtapu; wpiszKadr(); },
    etykiety: function (a) { if (a && a.length) etykiety = a.slice(); },
    tempo: function (v) { tempo = ogranicz(Number(v) || 1, 0.2, 4); },
    wstrzymaj: function () { biegnie = false; if (klatkaId) cancelAnimationFrame(klatkaId); klatkaId = 0; },
    wznow: function () {
      if (biegnie) return;
      biegnie = true; ostatniTS = 0;
      klatkaId = requestAnimationFrame(klatka);
    },
    zdejmij: function () {
      api.wstrzymaj();
      obserwatorMotywu.disconnect();
      if (obserwatorRozmiaru) obserwatorRozmiaru.disconnect(); else window.removeEventListener('resize', wymiar);
      if (obserwatorWidoku) obserwatorWidoku.disconnect();
      document.removeEventListener('visibilitychange', naWidocznosc);
      if (pytanieMotywu.removeEventListener) {
        pytanieMotywu.removeEventListener('change', przemaluj);
        pytanieRuchu.removeEventListener('change', naZmianeRuchu);
      }
      if (plotno.parentNode) plotno.parentNode.removeChild(plotno);
      delete host.bryla;
    }
  };

  host.bryla = api;
  wymiar();
  klatkaId = requestAnimationFrame(klatka);
  return api;
}

/* Zakładanie bryły na wszystkich polach dokumentu korzysta ze wspólnego mechanizmu przeglądania pól składnika, dostarczanego przez warstwę narzędzi wspólnych dla wszystkich składników okien. */
var zalozWszystkie = window.DanacoNarzedzia.polaSkladnika('data-bryla', zaloz);

window.DanacoBryla = { zaloz: zaloz, zalozWszystkie: zalozWszystkie };
})();
