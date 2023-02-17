package blend

var convertPy = `import bpy
import os
import shutil
from math import pi

from bpy.types import ViewLayer

class NpcType:
    def __init__(self, npcid, npcname):
        self.id = npcid
        self.name = npcname
npcs = []

class SpawnGroup:
    def __init__(self, id, spawngroupid):
        self.id = 0
        self.name = ""
        self.spawngroupid = 0
        self.spawn_limit = 0
        self.dist = 0
        self.max_x = 0
        self.min_x = 0
        self.max_y = 0
        self.min_y = 0
        self.delay = 0
        self.mindelay = 15000
        self.despawn = 0
        self.despawn_timer = 0
        self.wp_spawns = 0
spawngroups = {}

class Writer:
    def __init__(self, path):
        self.isCreated = False
        self.path = path
    def IsCreated(self):
        return self.isCreated
    def write(self, text):
        if not self.isCreated:
            if self.path.find("sql") and not os.path.exists(sql_path):
                os.makedirs(sql_path)
            self.w = open(self.path, "w+")
            self.isCreated = True
            if self.path.endswith("_EnvironmentEmitters.txt"):
                self.w.write("Name^EmitterDefIdx^X^Y^Z^Lifespan\n")
            if self.path.endswith("_doors.sql"):
                self.w.write("DELETE FROM doors WHERE zone = '"+base_name+"';\n")
                self.w.write("INSERT INTO doors (doorid, zone, ` + "`name`" + ` pos_x, pos_y, pos_z, heading, opentype, guild, lockpick, keyitem, nokeyring, triggerdoor, triggertype, disable_timer, doorisopen, door_param, dest_zone, dest_instance, dest_x, dest_y, dest_z, dest_heading, invert_state, incline, size, buffer, client_version_mask, is_ldon_door, min_expansion, max_expansion) VALUES\n")
            if self.path.endswith("_object.sql"):
                self.w.write("DELETE FROM object WHERE zoneid = "+zoneid+";\n")
                self.w.write("INSERT INTO object (zoneid, ` + "`version`" + `, xpos, ypos, zpos, heading, itemid, charges, objectname, ` + "`type`" + `, icon, unknown08, unknown10, unknown20, unknown24, unknown60, unknown64, unknown68, unknown72, unknown76, unknown84, size, tilt_x, tilt_y, min_expansion, max_expansion) VALUES\n")
        self.w.write(text)

def eulerToHeading(value):
    return round(180/pi*value/360*512)

def roundFloatStr(value):
    return str(round(value, 4))


print("eqgzi v1.8.0 converter")

blend_file_path = bpy.data.filepath
directory = os.path.dirname(blend_file_path)

base_name = os.path.basename(blend_file_path)
if base_name.find(".blend") != -1:
    base_name = base_name.replace(".blend", "")
out_path = directory + "/out"
cache_path = directory + "/cache"
sql_path = directory + "/sql"
zoneid = "32"

fe = Writer(out_path + "/" + base_name + "_EnvironmentEmitters.txt")
fsnd = Writer(out_path + "/" + base_name + ".emt")
fl = Writer(cache_path + "/" + base_name + "_light.txt")
fr = Writer(cache_path + "/" + base_name + "_region.txt")
fm = Writer(cache_path + "/" + base_name + "_material.txt")
fmod = Writer(cache_path + "/" + base_name + "_mod.txt")
fsg = Writer(sql_path + "/" + base_name + "_spawngroup_sql")
fs2 = Writer(sql_path + "/" + base_name + "_spawn2.sql")
fdoor = Writer(cache_path + "/" + base_name + "_doors.txt")
fdoorsql = Writer(sql_path + "/" + base_name + "_doors.sql")
fobjectsql = Writer(sql_path + "/" + base_name + "_object.sql")


print("Step 1) Deleting cache / out paths...")
# Delete contents of out path
if not os.path.exists(out_path):
    os.makedirs(out_path)
print("out path: " + out_path)


# Delete contents of cache path
if not os.path.exists(cache_path):
    os.makedirs(cache_path)
print("cache path: " + cache_path)

modDefs = {}

def isImageFile(name):
    for ext in {".dds", ".png", ".jpg"}:
        if name.find(ext) != -1:
            return True
    return False

def process(name, location, o):
    # check for any emitter definitions, any object can contain them
    if o.get("emit_id", 0) != 0:
        print("writing out emit_id "+str(o.get("emit_id", "1"))+" from object "+ name)
        fe.write(name + "^" + str(o.get("emit_id", "1")) + "^" + roundFloatStr(-location.y*2) + "^" + roundFloatStr(location.x*2) +"^" + roundFloatStr(location.z*2) + "^" + o.get("emit_duration", "90000000") + "\n")   
    if o.get("sound", 0) != 0:
        print("writing out sound "+str(o.get("sound", "1"))+" from object "+ name)
        fsnd.write("2,"+o.get("sound", "none.wav")+",0,")
        fsnd.write(str(o.get("sound_active", "0"))+",")
        fsnd.write(roundFloatStr(o.get("sound_volume", 1.0))+",")
        fsnd.write(str(o.get("sound_fade_in", "0"))+",")
        fsnd.write(str(o.get("sound_fade_out", "0"))+",")
        fsnd.write(str(o.get("sound_type", "0"))+",")
        fsnd.write(roundFloatStr(-location.y*2)+ ",")
        fsnd.write(roundFloatStr(location.x*2)+",")
        fsnd.write(roundFloatStr(location.z*2)+",")
        fsnd.write(roundFloatStr(o.get("sound_radius", 15.0))+",")
        fsnd.write(roundFloatStr(o.get("sound_distance", 50.0))+",")
        fsnd.write(str(o.get("sound_rand_distance", "0"))+",")
        fsnd.write(roundFloatStr(o.get("sound_trigger_range", 50.0))+",")
        fsnd.write(str(o.get("sound_min_repeat_delay", "0"))+",")
        fsnd.write(str(o.get("sound_max_repeat_delay", "0"))+",")
        fsnd.write(str(o.get("sound_max_repeat_delay", "0"))+",")
        fsnd.write(str(o.get("sound_xmi_index", "0"))+",")
        fsnd.write(str(o.get("sound_echo", "0"))+",")
        fsnd.write(str(o.get("sound_env_toggle", "1"))+"\n")
    if o.get("object_objectname", 0) != 0:
        print("writing out object "+str(o.get("object_objectname", "0"))+" from object "+name)
        if fobjectsql.IsCreated():
            fobjectsql.write(", \n")
        fobjectsql.write("("+str(o.get("object_zoneid", zoneid))+", ")
        fobjectsql.write(str(o.get("object_version", "0"))+", ")
        fobjectsql.write(roundFloatStr(-o.location.y*2) + ", " + roundFloatStr(o.location.x*2) + ", " + roundFloatStr(o.location.z*2)+", ")
        fobjectsql.write(str(eulerToHeading(o.rotation_euler.z))+", ") # heading
        fobjectsql.write(str(o.get("object_itemid", "0"))+", ")
        fobjectsql.write(str(o.get("object_charges", "0"))+", ")
        fobjectsql.write("'"+str(o.get("object_objectname", "0"))+"', ")
        fobjectsql.write(str(o.get("object_type", "0"))+", ")
        fobjectsql.write(str(o.get("object_icon", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown08", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown10", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown20", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown24", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown60", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown64", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown68", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown72", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown76", "0"))+", ")
        fobjectsql.write(str(o.get("object_unknown84", "0"))+", ")
        fobjectsql.write(str(o.get("object_size", "100"))+", ")
        fobjectsql.write(str(o.get("object_tilt_x", "0"))+", ")
        fobjectsql.write(str(o.get("object_tilt_y", "0"))+", ")
        fobjectsql.write(str(o.get("object_min_expansion", "0"))+", ")
        fobjectsql.write(str(o.get("object_max_expansion", "0"))+")")
    
    if o.get("spawngroup_id", "0") != "0":
        id = o.get("spawn2_id", 0)
        spawngroupid = o.get("spawngroupid", 0)
        if not spawngroupid in spawngroups:
            spawngroups[spawngroupid] = SpawnGroup(id, spawngroupid)
        if o.get("spawngroup_name", "0") != "0":
            spawngroups[spawngroupid].name = o["spawngroup_name"]
        if o.get("spawngroup_spawn_limit", "0") != "0":
            spawngroups[spawngroupid].spawn_limit = o["spawngroup_spawn_limit"]
        if o.get("spawngroup_dist", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_dist"]
        if o.get("spawngroup_max_x", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_max_x"]
        if o.get("spawngroup_max_y", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_max_y"]
        if o.get("spawngroup_min_x", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_min_x"]
        if o.get("spawngroup_min_y", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_min_y"]
        if o.get("spawngroup_delay", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_delay"]
        if o.get("spawngroup_mindelay", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_mindelay"]
        if o.get("spawngroup_despawn", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_despawn"]
        if o.get("spawngroup_despawn_timer", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_despawn_timer"]
        if o.get("spawngroup_wp_spawns", "0") != "0":
            spawngroups[spawngroupid].dist = o["spawngroup_wp_spawns"]
        fs2.write("REPLACE INTO spawn2 (id, spawngroupid, x, y, z, heading, respawntime, variance, pathgrid, version) VALUES(")
        fs2.write(str(o.get("spawn2_id", "0")) + ", " +str(o.get("spawn2_spawngroupid", "0"))+ ", ")
        fs2.write(str(location.x*2)+", "+str(location.y*2)+", "+str(location.z*2)+", ")
        fs2.write(str(eulerToHeading(o.rotation_euler.z))+ ", "+str(o.get("spawn2_respawntime", "0"))+ ", ")
        fs2.write(str(o.get("spawn2_variance", "0"))+ ", "+str(o.get("spawn2_pathgrid", "0"))+", ")
        fs2.write(str(o.get("spawn2_version", "0"))+");\n")

    if o.type == 'LIGHT':
        li = o.data
        if li.type == 'POINT':
            lightName = name.replace(" ", "-")
            if not lightName.startswith("LIB_") and not lightName.startswith("LIT_"):
                lightName = "LIB_" + lightName
            fl.write(lightName + " " + roundFloatStr(location.x*2) + " " + roundFloatStr(-location.y*2) + " " + roundFloatStr(location.z*2) + " " + roundFloatStr(li.color[0]) + " " + roundFloatStr(li.color[2]) + " " + roundFloatStr(li.color[1]) + " " + roundFloatStr(li.energy/10) + "\n")

    if o.type == 'EMPTY':
        if o.empty_display_type == 'CUBE':
            fr.write(name.replace(" ", "-") + " " + roundFloatStr(-location.y*2) + " " + roundFloatStr(location.x*2) + " " + roundFloatStr(location.z*2) + " " + roundFloatStr(o.scale.y*2) + " " + roundFloatStr(-o.scale.x*2) + " " + roundFloatStr((o.scale.z)*2) + " " + roundFloatStr(o.get("unknowna", 0)) + " " + roundFloatStr(o.get("unknownb", 0)) + " " + roundFloatStr(o.get("unknownc", 0)) + "\n")

print("Step 2) Applying modifiers...")
for o in bpy.data.objects:
    if not o.visible_get(view_layer=bpy.context.view_layer):
        continue
    bpy.context.view_layer.objects.active = o
    for mod in o.modifiers:
        print("applying modifier " + mod.name + " for " + o.name)
        bpy.ops.object.modifier_apply(modifier=mod.name)
bpy.ops.object.mode_set(mode = 'OBJECT')

print("Step 3) Writing material properties...")
for m in bpy.data.materials:
    fm.write("m " + m.name.replace(" ", "-") + " " + str(m.get("flag", 65536)) + " " + str(m.get("fx", "Opaque_MaxCB1.fx")) + "\n")
    for tree in m.texture_paint_slots:
        print("tree " +tree)
        for node in tree:
            print(node.name)
    try:
        name = m.node_tree.nodes['Image Texture'].image.name
    except:
        continue
    if name.find("."):
        name = name[0:name.find(".")]
    if os.path.isfile(name + ".txt"):
        mats = open(name + ".txt")
        lines = mats.readlines()
        for line in lines:
            if not isImageFile(line):
                continue
            line = line.replace("\n", "")
            if not os.path.isfile(directory+"/"+line):
                print(name+".txt: could not find animated texture " +directory+"/"+line)
                exit(1)
                continue
            print(name+".txt: copying animated texture to cache: "+line)
            shutil.copyfile(line, "cache/"+line)
        shutil.copyfile(name+".txt", "cache/"+name+".txt")
        
            
    for prop in m.items():
        for k in prop:
            if not isinstance(k, str):
                continue
            if not k.startswith("e_"):
                continue
            eValue = str(m[k])
            if eValue.find(" ") == -1:
                eValue = "0 "+eValue                
            fm.write("e " + m.name.replace(" ", "-") + " " +  k + " " + eValue +"\n")
            for entry in eValue.split(" "):
                if not isImageFile(entry):
                    continue
                if not os.path.exists(entry):
                    print("failed to find "+entry+" in current path, defined on material "+m.name)
                    exit(1)
                print("copying "+entry+" to cache")
                shutil.copyfile(entry, "cache/"+entry)



print("Step 4) Removing any hidden objects...")
for o in bpy.data.objects:
    if not o.visible_get(view_layer=bpy.context.view_layer): 
        print("removing " + o.name + " (not active view)")
        bpy.data.objects.remove(o, do_unlink=True)
        continue


print("Step 5) Exporting any linked objects...")
exportedMods = []

for o in bpy.data.objects:
    if not o.visible_get(view_layer=bpy.context.view_layer):
        print(o.name + "skipped, it is not visible for export")
        continue
    col = o.instance_collection
    if not col:
        continue
    print(col.name +" has a link instance as "+o.name)
    
    col.library.reload()
    col = o.instance_collection
    for co in col.objects:
        print(co.name+ " found and processed")
        process(o.name, o.location+co.location, co)
        if co.type != 'MESH':
            bpy.data.objects.remove(co, do_unlink=True)
    if not col.library:
        print(col.name +" has no library data, skipping export")
        continue
    bpy.ops.object.select_all(action='DESELECT')
    isExported = False
    for e in exportedMods:
        if not e == col.name:
            continue
        isExported = True
        break
    if not isExported:
        print(col.name + " is going to be exported from " +col.library.name)
        exportedMods.append(col.name)

    objName = col.library.name.replace(".blend", ".obj")

    if o.get("door_id", "0") != "0":
        print(col.name + " has door data")
        fdoor.write(objName + "\n")
        if fdoorsql.IsCreated():
            fdoorsql.write(", \n")
        fdoorsql.write("(" + str(o.get("door_id", "0"))+", ")
        fdoorsql.write("'"+base_name+"', ")
        fdoorsql.write("'"+objName.replace(" ", "-").replace(".obj", "").upper()+"', ")
        fdoorsql.write(roundFloatStr(o.location.x*2) + ", " + roundFloatStr(-o.location.y*2) + ", " + roundFloatStr(o.location.z*2)+", ")
        
        fdoorsql.write(str(eulerToHeading(o.rotation_euler.z))+", ") # heading
        fdoorsql.write(str(o.get("door_opentype", "0"))+", ")
        fdoorsql.write(str(o.get("door_guild", "0"))+", ")
        fdoorsql.write(str(o.get("door_lockpick", "0"))+", ")
        fdoorsql.write(str(o.get("door_keyitem", "0"))+", ")
        fdoorsql.write(str(o.get("door_nokeyring", "0"))+", ")
        fdoorsql.write(str(o.get("door_door_triggerdoor", "0"))+", ")
        fdoorsql.write(str(o.get("door_door_triggertype", "0"))+", ")
        fdoorsql.write(str(o.get("door_disable_timer", "0"))+", ")
        fdoorsql.write(str(o.get("door_doorisopen", "0"))+", ")
        fdoorsql.write(str(o.get("door_param", "0"))+", ")
        fdoorsql.write("'"+str(o.get("door_dest_zone", "NONE"))+"', ")
        fdoorsql.write(str(o.get("door_dest_instance", "0"))+", ")
        fdoorsql.write(str(o.get("door_dest_x", "0"))+", ")
        fdoorsql.write(str(o.get("door_dest_y", "0"))+", ")
        fdoorsql.write(str(o.get("door_dest_z", "0"))+", ")
        fdoorsql.write(str(o.get("door_dest_heading", "0"))+", ")
        fdoorsql.write(str(o.get("door_invert_state", "0"))+", ")
        fdoorsql.write(str(o.get("door_incline", "0"))+", ")
        fdoorsql.write(str(o.get("door_size", "100"))+", ")
        fdoorsql.write(str(o.get("door_buffer", "0"))+", ")
        fdoorsql.write(str(o.get("door_client_version_mask", "4294967295"))+", ")
        fdoorsql.write(str(o.get("door_is_ldon_door", "0"))+", ")
        fdoorsql.write(str(o.get("door_min_expansion", "0"))+", ")
        fdoorsql.write(str(o.get("door_max_expansion", "0"))+")")

    else:
        fmod.write(objName + " " + o.name.replace(" ", "-") + " " + roundFloatStr(-o.location.y*2) + " " + roundFloatStr(o.location.x*2) + " " + roundFloatStr(o.location.z*2) + " "  + roundFloatStr(-o.rotation_euler.y) + " " + roundFloatStr(o.rotation_euler.x) + " " + roundFloatStr(o.rotation_euler.z) + " " + roundFloatStr(o.scale.z) + "\n")
    if isExported:
        print(col.name+" is already exported, only adding placement instance data")
        for co in col.objects:
            bpy.data.objects.remove(co, do_unlink=True)
        continue
    obj_file = os.path.join(cache_path, objName)
    
    #for attr in dir(col):
    #    print("col.%s = %r" % (attr, getattr(col, attr)))
    for co in col.objects:
        bpy.context.scene.collection.objects.link(co)
        bpy.context.view_layer.objects.active = co
        co.select_set(True)
    bpy.ops.export_scene.obj(filepath=obj_file, check_existing=True, axis_forward='-X', axis_up='Z', filter_glob="*.obj;*.mtl", use_selection=True, use_animation=False, use_mesh_modifiers=True, use_edges=True, use_smooth_groups=False, use_smooth_groups_bitflags=False, use_normals=True, use_uvs=True, use_materials=True, use_triangles=True, use_nurbs=False, use_vertex_groups=False, use_blen_objects=True, group_by_object=False, group_by_material=False, keep_vertex_order=False, global_scale=2, path_mode='COPY')
    bpy.data.objects.remove(o, do_unlink=True)


if fdoorsql.IsCreated():
    fdoorsql.write(";\n")

print("Step 6) Processing zone objects...")
bpy.ops.object.select_all(action='DESELECT')
for o in bpy.data.objects:
    print(o.name + " processing (" + o.type + ")")
    process(o.name, o.location, o)
    if o.type != 'MESH':
        bpy.data.objects.remove(o, do_unlink=True)
        continue


if fobjectsql.IsCreated():
    fobjectsql.write(";\n")

for sp in spawngroups:
    fsg.write("REPLACE INTO spawngroup (id, name, spawn_limit, dist, max_x, min_x, max_x, min_y, delay, mindelay, despawn, despawn_timer, wp_spawns) VALUES ("+str(spawngroups[sp].id)+", '"+str(spawngroups[sp].name)+"', "+str(spawngroups[sp].spawn_limit)+", "+str(spawngroups[sp].dist)+", "+str(spawngroups[sp].max_x)+", "+str(spawngroups[sp].min_x)+", "+str(spawngroups[sp].max_x)+", "+str(spawngroups[sp].min_y)+", "+str(spawngroups[sp].delay)+", "+str(spawngroups[sp].mindelay)+", "+str(spawngroups[sp].despawn)+", "+str(spawngroups[sp].despawn_timer)+", "+str(spawngroups[sp].wp_spawns)+");\n")

print("Step 7) Exporting zone .obj")
bpy.ops.export_scene.obj(filepath=cache_path + "/" + base_name + '.obj', check_existing=True, axis_forward='-X', axis_up='Z', filter_glob="*.obj;*.mtl", use_selection=False, use_animation=False, use_mesh_modifiers=True, use_edges=True, use_smooth_groups=False, use_smooth_groups_bitflags=False, use_normals=True, use_uvs=True, use_materials=True, use_triangles=True, use_nurbs=False, use_vertex_groups=False, use_blen_objects=True, group_by_object=False, group_by_material=False, keep_vertex_order=False, global_scale=2, path_mode='COPY')
exit(0)`
