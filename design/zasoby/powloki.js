/* Powłoki platformy nadlatują, układają się w wachlarz i scalają w jeden podświetlony blok ze znakiem marki, sygnalizując zakończenie montażu platformy. ================
   POWŁOKI PLATFORMY — składnik biblioteki

   Cztery powłoki programu — rdzeń, środowiska pracy, moduły, stanowisko —
   nadlatują, układają się w wachlarz i scalają w jeden blok ze znakiem marki
   na wierzchu. Rzecz do lewej kolumny okna „Przygotowanie środowiska pracy":
   pokazuje, co właśnie się składa, zamiast zajmować oko samym ruchem.

   Pętla 3,00 s:
       0,00–1,20  powłoki nadlatują i układają się w wachlarz
       1,20–1,90  wachlarz scala się w blok, krawędzie zapalają się
       1,90–2,78  blok obraca się, po blatach przechodzi refleks
       2,78–3,00  wygaszenie i start od nowa

   Okno wstawia puste pole i nic więcej:

       <div class="dn-powloki" data-powloki
            role="img" aria-label="Cztery powłoki platformy Danaco Console"></div>

   Składnik NIE MA własnego tła — kładzie się wprost na powierzchni okna.
   Barwy wyłącznie z żetonów `--dn-bryla-*`, tych samych, którymi rządzi się
   bryła instalatora: jedna paleta platformy, nie dwie podobne.

   Sterowanie z okna:
       var p = document.querySelector('[data-powloki]').powloki;
       p.setProgress(0.42);   // 0..1 — składanie idzie za postępem
       p.loop();              // powrót do pętli automatycznej
       p.setLabels(['RDZEŃ', 'ŚRODOWISKA', 'MODUŁY', 'STANOWISKO']);
       p.setSpeed(1.2); p.pause(); p.resume(); p.destroy();

   Atrybuty pola:
       data-cykl="3"       długość pętli w sekundach
       data-podpisy="0"    bez podpisów na powłokach
       data-znak="0"       bez znaku marki na wierzchu
       data-tempo="1"      mnożnik tempa
       data-postep="0.5"   postęp startowy — wyłącza pętlę
   ============================================================================ */
(function () {
'use strict';

/* Czytnik barw stoi w `narzedzia-okien.js` — potrzebuje go każdy składnik
   rysujący na płótnie. */
var czytnikBarw = window.DanacoNarzedzia.czytnikBarw;

function zaloz(host) {
  if (host.powloki) return host.powloki;
  /* Pole bywa kopią znacznika, w której stoi już płótno bez wiązania — powstaje
     tak wszędzie, gdzie okno składa się z klonu szablonu. Zdejmujemy je, żeby
     nie rysować na martwym płótnie pod nowym. */
  var stare = host.querySelectorAll('canvas');
  for (var s = 0; s < stare.length; s++) stare[s].parentNode.removeChild(stare[s]);

  var st=host, cv=document.createElement('canvas'), ctx=cv.getContext('2d');
  host.appendChild(cv);
  var dane=host.dataset;
  var REDUCE=window.matchMedia&&window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  function clamp(v,a,b){return v<a?a:(v>b?b:v);}
  function ss(t,a,b){return clamp((t-a)/(b-a),0,1);}
  function eOut(t){return 1-Math.pow(1-clamp(t,0,1),3);}
  function eInOut(t){t=clamp(t,0,1);return t<.5?4*t*t*t:1-Math.pow(-2*t+2,3)/2;}
  function lerp(a,b,t){return a+(b-a)*t;}
  function rgba(c,a){return 'rgba('+c[0]+','+c[1]+','+c[2]+','+a+')';}

  /* Paleta z żetonów — tych samych, którymi rządzi się bryła instalatora.
     Jedna paleta platformy, nie dwie podobne: rysunek dostarczony osobno miał
     blaty o włos ciemniejsze, a różnica nie niosła żadnego rozstrzygnięcia. */
  var cz=czytnikBarw(host);
  var P={
    line:   cz.barwa('--dn-bryla-linia',    '#68A8F6').rgb,
    glow:   cz.barwa('--dn-bryla-poswiata', '#3474C8').rgb,
    soft:   cz.barwa('--dn-bryla-lagodna',  '#84B6F0').rgb,
    node:   cz.barwa('--dn-bryla-wezel',    '#A8CEF6').rgb,
    hot:    cz.barwa('--dn-bryla-blysk',    '#F0F8FF').rgb,
    accent: cz.barwa('--dn-znak-poswiata',  '#5C8CEC').rgb,
    label:  cz.barwa('--dn-bryla-etykieta', '#96B2D6').rgb,
    top:    cz.barwa('--dn-bryla-blat',         'rgba(41,66,101,.97)').css,
    topHi:  cz.barwa('--dn-bryla-blat-wierzch', 'rgba(52,82,122,.97)').css,
    sideA:  cz.barwa('--dn-bryla-bok',          'rgba(26,45,71,.97)').css,
    sideB:  cz.barwa('--dn-bryla-bok-cien',     'rgba(17,30,48,.97)').css,
    shadow: cz.barwa('--dn-bryla-cien',         'rgba(16,36,66,.26)').css
  };
  var ZNAK_LICO = cz.barwa('--dn-znak-lico', '#F4F4F4').css;
  cz.zdejmij();

  var W=0,H=0,DPR=1,unit=1,cX=0,cY=0;
  var cam={yaw:.62,pitch:.50,dist:6.1,f:3.0};
  function p3(x,y,z){
    var cy=Math.cos(cam.yaw),sy=Math.sin(cam.yaw);
    var rx=x*cy-z*sy, rz=x*sy+z*cy;
    var cp=Math.cos(cam.pitch),sp=Math.sin(cam.pitch);
    var yc=y*cp+rz*sp, zc=-y*sp+rz*cp+cam.dist;
    var s=cam.f*unit/Math.max(.08,zc);
    return {x:cX+rx*s,y:cY-yc*s,s:s,d:zc};
  }
  function poly(p){ ctx.beginPath(); ctx.moveTo(p[0].x,p[0].y);
    for(var i=1;i<p.length;i++) ctx.lineTo(p[i].x,p[i].y); ctx.closePath(); }

  /* ------------------------------------------------ powłoki */
  var PW=1.34, PD=0.94, TH=0.05;             // szerokość, głębokość, grubość
  var FANX=0.16, FANZ=0.20, GAP0=0.30, GAP1=0.115;
  var LB=['STANOWISKO','MODUŁY','ŚRODOWISKA PRACY','RDZEŃ PLATFORMY'];
  var showLabels=dane.podpisy!=='0', showMark=dane.znak!=='0';
  var shells=[];
  function build(){
    shells=[];
    for(var i=0;i<4;i++){                       // i=0 wierzch ... 3 spód
      var dir=(i%2?1:-1);
      shells.push({
        i:i, top:(i===0),
        from:{x:dir*(2.4+i*0.25), y:1.5-i*0.55, z:(i<2?-1:1)*(1.5+i*0.2)},
        rotFrom:dir*0.55, a:0, flash:0
      });
    }
  }
  build();

  function tf(ax,ay,az,rot){                     // baza płaszczyzny blatu -> transformacja 2D
    var o=p3(ax,ay,az);
    var cs=Math.cos(rot), sn=Math.sin(rot);
    var ex=p3(ax+cs,ay,az+sn), ez=p3(ax-sn,ay,az+cs);
    return {a:ex.x-o.x,b:ex.y-o.y,c:ez.x-o.x,d:ez.y-o.y,e:o.x,f:o.y,s:o.s};
  }
  function planePt(ax,ay,az,rot,lx,lz){
    var cs=Math.cos(rot), sn=Math.sin(rot);
    return p3(ax+lx*cs-lz*sn, ay, az+lx*sn+lz*cs);
  }
  var CF=[[0,1,2,3],[4,5,6,7],[0,1,5,4],[1,2,6,5],[2,3,7,6],[3,0,4,7]];
  function shellVerts(ax,ay,az,rot,sc){
    var hw=PW/2*sc, hd=PD/2*sc, hh=TH/2, v=[];
    var sg=[[-1,1,-1],[1,1,-1],[1,1,1],[-1,1,1],[-1,-1,-1],[1,-1,-1],[1,-1,1],[-1,-1,1]];
    var cs=Math.cos(rot), sn=Math.sin(rot);
    for(var i=0;i<8;i++){
      var x=sg[i][0]*hw, y=sg[i][1]*hh, z=sg[i][2]*hd;
      v.push(p3(ax+(x*cs-z*sn), ay+y, az+(x*sn+z*cs)));
    }
    return v;
  }
  var sygPaths=null;
  try{ sygPaths=[new Path2D('M12 26 H24 L44 48.0 L24 70 H12 L32 48.0 Z'),
                new Path2D('M40 26 H52 L72 48.0 L52 70 H40 L60 48.0 Z')]; }catch(e){}

  function drawShell(S,alpha,tot,merge,sweep){
    var k=eOut(S.a);
    var fan=1-merge;
    var bx=(1.5-S.i)*FANX*fan, bz=(S.i-1.5)*FANZ*fan;
    var gap=lerp(GAP0,GAP1,merge);
    var by=(1.5-S.i)*gap;
    var ax=lerp(S.from.x,bx,k), ay=lerp(S.from.y,by,k), az=lerp(S.from.z,bz,k);
    var rot=lerp(S.rotFrom,0,k);
    var sc=lerp(0.9,1,k);
    var v=shellVerts(ax,ay,az,rot,sc), i, f=[];
    for(i=0;i<CF.length;i++){ var q=CF[i]; f.push({q:q,i:i,d:(v[q[0]].d+v[q[1]].d+v[q[2]].d+v[q[3]].d)/4}); }
    f.sort(function(a,b){return b.d-a.d;});
    ctx.save(); ctx.globalAlpha=alpha;
    for(i=0;i<f.length;i++){
      var F=f[i], pts=[v[F.q[0]],v[F.q[1]],v[F.q[2]],v[F.q[3]]];
      poly(pts);
      ctx.fillStyle=F.i===0?(S.top?P.topHi:P.top):(F.i===1?P.sideB:P.sideA);
      ctx.fill();
      ctx.strokeStyle=rgba(P.line,(F.i===0?.85:.45)+.5*S.flash);
      ctx.lineWidth=(F.i===0?1.2:1)+1.8*S.flash;
      ctx.stroke();
    }
    // zawartość blatu: siatka, podpis, znak
    if(k>0.55){
      var ca=ss(k,0.55,0.92);
      ctx.save();
      poly([v[0],v[1],v[2],v[3]]); ctx.clip();
      var T=tf(ax,ay+TH/2+0.001,az,rot);
      ctx.transform(T.a,T.b,T.c,T.d,T.e,T.f);      // od tej chwili rysujemy w układzie blatu
      // siatka
      ctx.strokeStyle=rgba(P.line,0.13*ca);
      ctx.lineWidth=0.006;
      var g;
      for(g=-PW/2;g<=PW/2+1e-6;g+=PW/9){ ctx.beginPath(); ctx.moveTo(g,-PD/2); ctx.lineTo(g,PD/2); ctx.stroke(); }
      for(g=-PD/2;g<=PD/2+1e-6;g+=PD/6){ ctx.beginPath(); ctx.moveTo(-PW/2,g); ctx.lineTo(PW/2,g); ctx.stroke(); }
      // refleks przesuwający się po blacie
      if(sweep>0){
        var sx=-PW/2+PW*sweep;
        var lg=ctx.createLinearGradient(sx-PW*0.16,0,sx+PW*0.16,0);
        lg.addColorStop(0,rgba(P.soft,0)); lg.addColorStop(.5,rgba(P.soft,.13)); lg.addColorStop(1,rgba(P.soft,0));
        ctx.fillStyle=lg; ctx.fillRect(-PW/2,-PD/2,PW,PD);
      }
      // podpis powłoki
      if(showLabels){
        ctx.save();
        ctx.scale(1,-1);
        ctx.font='0.088px "JetBrains Mono",ui-monospace,monospace';
        ctx.textBaseline='middle';
        ctx.fillStyle=rgba(P.label,0.85*ca);
        ctx.fillText(LB[S.i], -PW/2+0.075, PD/2-0.085);
        ctx.restore();
      }
      // znak marki na wierzchniej powłoce
      if(S.top&&showMark&&sygPaths){
        var ma=ss(k,0.8,1)*ca;
        ctx.save();
        ctx.translate(-0.06,0.02); ctx.scale(1,-1);
        var msc=0.0072;
        ctx.scale(msc,msc); ctx.translate(-48,-48);
        ctx.globalAlpha=ma;
        ctx.fillStyle=ZNAK_LICO;
        ctx.fill(sygPaths[0]); ctx.fill(sygPaths[1]);
        ctx.fillStyle=rgba(P.accent,1);
        ctx.beginPath(); ctx.arc(83,63.5,6.5,0,Math.PI*2); ctx.fill();
        ctx.restore();
      }
      ctx.restore();
    }
    if(S.flash>0.02){ poly([v[0],v[1],v[2],v[3]]); ctx.fillStyle=rgba(P.hot,.26*S.flash); ctx.fill(); }
    ctx.restore();
    return {x:ax,y:ay,z:az};
  }

  /* ------------------------------------------------ kadr */
  function fitPts(){
    var a=[],i,R=Math.hypot(PW/2+FANX,PD/2+FANZ)*1.02, yT=GAP0*1.5+0.16, yB=-GAP0*1.5-0.16;
    for(i=0;i<20;i++){ var t=i/20*Math.PI*2;
      a.push([Math.cos(t)*R,yB,Math.sin(t)*R]);
      a.push([Math.cos(t)*R,yT,Math.sin(t)*R]);
    }
    return a;
  }
  function resize(){
    DPR=Math.min(2,window.devicePixelRatio||1);
    W=st.clientWidth||1; H=st.clientHeight||1;
    cv.width=Math.max(1,Math.round(W*DPR)); cv.height=Math.max(1,Math.round(H*DPR));
    cv.style.width=W+'px'; cv.style.height=H+'px';
    ctx.setTransform(DPR,0,0,DPR,0,0);
    unit=Math.min(W,H*1.3)*0.42; cX=W*0.5; cY=H*0.5;
    fit();
  }
  function fit(){
    var yaw0=cam.yaw,pitch0=cam.pitch;
    var pad=Math.max(6,Math.min(W,H)*0.055);
    var bw=W-2*pad, bh=H-2*pad, pts=fitPts(), i;
    for(var it=0;it<3;it++){
      var ys=[0.28,0.50,0.72,0.95,1.10], ps=[0.44,0.50,0.56];
      var minx=1e9,maxx=-1e9,miny=1e9,maxy=-1e9,yi,pi;
      for(yi=0;yi<ys.length;yi++){ cam.yaw=ys[yi];
        for(pi=0;pi<ps.length;pi++){ cam.pitch=ps[pi];
          for(i=0;i<pts.length;i++){
            var q=p3(pts[i][0],pts[i][1],pts[i][2]);
            if(q.x<minx)minx=q.x; if(q.x>maxx)maxx=q.x;
            if(q.y<miny)miny=q.y; if(q.y>maxy)maxy=q.y;
          }
        }
      }
      cam.yaw=yaw0; cam.pitch=pitch0;
      var k=Math.min(bw/Math.max(1,maxx-minx),bh/Math.max(1,maxy-miny));
      k=Math.max(0.3,Math.min(1.7,k)); unit*=k;
      var mx=(minx+maxx)/2,my=(miny+maxy)/2;
      cX+=(pad+bw/2)-(cX+(mx-cX)*k);
      cY+=(pad+bh/2)-(cY+(my-cY)*k);
    }
  }

  /* ------------------------------------------------ cykl 3 s */
  var TOTAL=clamp(parseFloat(dane.cykl||'3')||3,1.5,12);
  var T={arr0:0.00,arr1:0.40,mrg0:0.40,mrg1:0.63,rot1:0.93,end:1.00};   // ułamki cyklu
  var t=0,running=true,last=0,tot=0,external=false,progress=0;
  var speed=parseFloat(dane.tempo||'1')||1;
  if(dane.postep!==undefined){ external=true; progress=clamp(parseFloat(dane.postep)||0,0,1); }

  function reset(){ t=0; for(var i=0;i<shells.length;i++){ shells[i].a=0; shells[i].flash=0; } }
  function update(dt){
    tot+=dt;
    if(!external){ t+=dt; if(t>=TOTAL) reset(); }
    var u=external?progress:(t/TOTAL);
    var arrive=ss(u,T.arr0,T.arr1);
    var merge=eInOut(ss(u,T.mrg0,T.mrg1));
    // obrót: powolny w fazie montażu, wyraźny po scaleniu
    var spin=eInOut(ss(u,T.mrg1,T.rot1));
    cam.yaw=0.30+arrive*0.14+spin*0.62+Math.sin(tot*0.5)*0.02;
    cam.pitch=0.50-0.05*merge+Math.sin(tot*0.4)*0.015;
    for(var i=0;i<shells.length;i++){
      var S=shells[i];
      var tgt=clamp((arrive*4-(3-i))/0.62,0,1), pr=S.a;   // od spodu do wierzchu
      S.a+=(tgt-S.a)*Math.min(1,dt*12);
      if(pr<0.93&&S.a>=0.93) S.flash=1;
      S.flash=Math.max(0,S.flash-dt*4);
    }
    return {u:u,merge:merge,spin:spin};
  }
  function render(state){
    ctx.clearRect(0,0,W,H);
    var u=state.u;
    var al = external?1:(1-ss(u,T.rot1,T.end));
    if(al<=0.01) return;
    var sc=p3(0,-GAP0*1.5-0.06,0);
    ctx.save(); ctx.globalAlpha=al*ss(u,0.05,0.5);
    var g=ctx.createRadialGradient(sc.x,sc.y,1,sc.x,sc.y,unit*0.66);
    g.addColorStop(0,P.shadow); g.addColorStop(1,'transparent');
    ctx.fillStyle=g; ctx.fillRect(0,0,W,H); ctx.restore();

    var sweep=state.merge>0.9 ? ((tot*0.55)%1) : 0;
    var list=shells.slice().sort(function(a,b){ return b.i-a.i; });   // spód najpierw
    for(var i=0;i<list.length;i++){
      var S=list[i];
      if(S.a<=0.01) continue;
      drawShell(S, al*clamp(S.a*1.4,0,1), tot, state.merge, sweep);
    }
  }
  function frame(ts){
    if(!running) return;
    if(!last) last=ts;
    var dt=Math.min(0.05,(ts-last)/1000)*speed; last=ts;
    var s=update(dt); render(s);
    requestAnimationFrame(frame);
  }
  var api={
    setProgress:function(p){ external=true; progress=clamp(Number(p)||0,0,1); return progress; },
    getProgress:function(){ return external?progress:t/TOTAL; },
    loop:function(){ external=false; reset(); },
    setLabels:function(a){ if(a&&a.length){ for(var i=0;i<4&&i<a.length;i++) LB[i]=String(a[i]); } },
    setSpeed:function(v){ speed=clamp(Number(v)||1,.2,4); },
    setCycle:function(v){ TOTAL=clamp(Number(v)||3,1.5,12); },
    pause:function(){ running=false; },
    resume:function(){ if(!running){ running=true; last=0; requestAnimationFrame(frame); } },
    destroy:function(){ running=false; window.removeEventListener('resize',resize);
      if(cv&&cv.parentNode) cv.parentNode.removeChild(cv); }
  };
  host.powloki=api;
  window.addEventListener('resize',resize);
  document.addEventListener('visibilitychange',function(){ if(document.hidden) api.pause(); else api.resume(); });
  resize();
  requestAnimationFrame(frame);
  return api;
}

/* Zakładanie na wszystkich polach oznaczonych atrybutem danych korzysta ze wspólnego mechanizmu biblioteki narzędzi, dzielonego przez wszystkie składniki platformy. */
var zalozWszystkie = window.DanacoNarzedzia.polaSkladnika('data-powloki', zaloz);

window.DanacoPowloki = { zaloz: zaloz, zalozWszystkie: zalozWszystkie };
})();
