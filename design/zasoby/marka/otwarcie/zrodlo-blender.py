import bpy, math, sys, os
from mathutils import Vector

motyw = sys.argv[sys.argv.index('--')+1] if '--' in sys.argv else 'ciemny'
BARWA_ZNAK = (0.898,0.898,0.898,1) if motyw=='ciemny' else (0.094,0.094,0.094,1)   # F4F4F4 / 181818
BARWA_KROPKA = (0.361,0.549,0.925,1) if motyw=='ciemny' else (0.231,0.435,0.878,1) # 5C8CEC / 3B6FE0
KL = 90

def m(x,y): return ((x-48)/48.0, (48-y)/48.0)
GROT = [(12,26),(24,26),(44,48),(24,70),(12,70),(32,48)]

sc = bpy.context.scene
bpy.ops.object.select_all(action='SELECT'); bpy.ops.object.delete()

def grot(dx, nazwa):
    verts=[m(x+dx,y)+(0,) for x,y in GROT]
    me=bpy.data.meshes.new(nazwa); me.from_pydata([Vector(v) for v in verts],[],[list(range(6))])
    me.update(); ob=bpy.data.objects.new(nazwa,me); sc.collection.objects.link(ob)
    bpy.context.view_layer.objects.active=ob; ob.select_set(True)
    bpy.ops.object.mode_set(mode='EDIT'); bpy.ops.mesh.select_all(action='SELECT')
    bpy.ops.mesh.extrude_region_move(TRANSFORM_OT_translate={"value":(0,0,0.030)})
    bpy.ops.object.mode_set(mode='OBJECT')
    bpy.ops.object.modifier_add(type='BEVEL'); ob.modifiers[-1].width=0.0035; ob.modifiers[-1].segments=2
    ob.select_set(False); return ob

g1=grot(0,'grot1'); g2=grot(28,'grot2')
bpy.ops.mesh.primitive_cylinder_add(radius=6.5/48.0, depth=0.030, vertices=64,
    location=(m(83,63.5)[0], m(83,63.5)[1], 0.015))
kropka=bpy.context.object; kropka.name='kropka'
bpy.ops.object.modifier_add(type='BEVEL'); kropka.modifiers[-1].width=0.0035; kropka.modifiers[-1].segments=2

def mat(nazwa,barwa,metal,rough):
    ma=bpy.data.materials.new(nazwa); ma.use_nodes=True
    b=ma.node_tree.nodes['Principled BSDF']
    b.inputs['Base Color'].default_value=barwa
    b.inputs['Metallic'].default_value=metal; b.inputs['Roughness'].default_value=rough
    return ma
mz=mat('znak',BARWA_ZNAK,0.0,0.62); mk=mat('kropka',BARWA_KROPKA,0.0,0.55)
for o in (g1,g2): o.data.materials.append(mz)
kropka.data.materials.append(mk)

pref=bpy.context.preferences.edit
pref.keyframe_new_interpolation_type='LINEAR'

# Ruch liczony klatka po klatce według własnej krzywej: rozkład narastania jest równomierny, funkcja nie zawiera odcinków o zerowej prędkości zmiany.
def ease(t):            # wyjscie czwartego stopnia: rusza od razu, osiada lagodnie
    return 1-(1-t)**3
def lerp(a,b,u): return a+(b-a)*u

def anim(ob, k0, k1, poz0, poz1, rot0, rot1):
    for f in range(1, KL+1):
        u = 0.0 if f<=k0 else (1.0 if f>=k1 else ease((f-k0)/(k1-k0)))
        ob.location=[lerp(poz0[i],poz1[i],u) for i in range(3)]
        ob.rotation_euler=[lerp(rot0[i],rot1[i],u) for i in range(3)]
        ob.keyframe_insert('location',frame=f); ob.keyframe_insert('rotation_euler',frame=f)

R=math.radians
anim(g1, 1, 80, (0,0,-0.14),(0,0,0), (0,R(38),0),(0,0,0))
anim(g2, 7, 86, (0,0,-0.14),(0,0,0), (0,R(38),0),(0,0,0))

for f in range(1, KL+1):                      # kropka: wejscie z lekkim przeskokiem
    u = 0.0 if f<=58 else (1.0 if f>=88 else ease((f-58)/30.0))
    sk = u*(1.06-0.06*u)
    kropka.scale=(max(sk,0.001),max(sk,0.001),1)
    kropka.keyframe_insert('scale',frame=f)

bpy.ops.object.camera_add(location=(0,0,3.9), rotation=(0,0,0))
cam=bpy.context.object; sc.camera=cam; cam.data.lens=68
anim(cam, 1, 88, (0.26,-0.30,3.80),(0.0,0.0,3.9), (R(4.5),0,R(3.5)),(0,0,0))

for loc,en in [((2.0,-1.6,5.0),300),((-2.4,-1.0,4.4),260),((0.2,2.4,4.0),220),((0,0,6.0),240)]:
    bpy.ops.object.light_add(type='AREA',location=loc); l=bpy.context.object
    l.data.energy=en; l.data.size=7.0
    d=Vector((0,0,0))-Vector(loc); l.rotation_euler=d.to_track_quat('-Z','Y').to_euler()

sc.render.engine='CYCLES'; sc.cycles.samples=64; sc.cycles.device='CPU'
sc.cycles.use_denoising=True
sc.render.resolution_x=460; sc.render.resolution_y=460
sc.render.use_motion_blur=True
sc.render.motion_blur_shutter=0.85
sc.render.fps=30
sc.render.film_transparent=False
TLO = (0.0075,0.0075,0.0075,1) if motyw=='ciemny' else (1.0,1.0,1.0,1)   # #181818 / #FFFFFF w przestrzeni liniowej
sw=bpy.data.worlds.new('tlo'); sc.world=sw; sw.use_nodes=True
sw.node_tree.nodes['Background'].inputs['Color'].default_value=TLO
sw.node_tree.nodes['Background'].inputs['Strength'].default_value=1.0
sc.render.image_settings.file_format='PNG'; sc.render.image_settings.color_mode='RGBA'
sc.frame_start=1; sc.frame_end=KL
sc.render.filepath=f'/tmp/claude-1000/-home-ubuntu/797585e4-0f17-4f1a-a125-239f501741e7/scratchpad/b3d/kl-{motyw}-'
bpy.ops.render.render(animation=True)
print('GOTOWE', motyw)
