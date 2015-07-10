package aefs

import "appengine/datastore"

// DeleteAll deletes across all roots
// DeleteAll deletes by kind alone.
func (fs *AeFileSys) DeleteAll() (string, error) {

	msg := ""
	{
		q := datastore.NewQuery(tfil).KeysOnly()
		var files []AeFile
		keys, err := q.GetAll(fs.Ctx(), &files)
		if err != nil {
			msg += "could not get file keys\n"
			return msg, err
		}

		err = datastore.DeleteMulti(fs.Ctx(), keys)
		if err != nil {
			msg += "error deleting files\n"
			return msg, err
		}

		msg += spf("%v files deleted\n", len(keys))

	}

	{
		q := datastore.NewQuery(tdir).KeysOnly()
		var dirs []AeDir
		keys, err := q.GetAll(fs.Ctx(), &dirs)
		if err != nil {
			msg += "could not get dir keys\n"
			return msg, err
		}

		err = datastore.DeleteMulti(fs.Ctx(), keys)
		if err != nil {
			msg += "error deleting directories\n"
			return msg, err
		}

		msg += spf("%v directories deleted\n", len(keys))
	}

	return msg, nil
}
